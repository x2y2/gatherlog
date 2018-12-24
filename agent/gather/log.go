package logs

import (
	"os"
	"io/ioutil"
	"strconv"
	"gatherlog/agent/common"
	"sync"
	"gatherlog/agent/chkpoint"
	"path/filepath"
	"strings"
	"log"
	"io"
)


//从文件尾开始读
func TailLog(offset int64,logfile string) string{
	fd,err := os.Open(logfile)
	if err != nil{
		log.Fatal(err)
	}
	defer fd.Close()

	//从offset处读取日志
	fd.Seek(offset,0)
	data , err := ioutil.ReadAll(fd)
	if err != nil{log.Fatal(err)}

	//当前offset写到chkpt 文件
	cur_offset,_ := fd.Seek(0,1)
	str_offset := strconv.Itoa(int(cur_offset))

	c := common.Config{}
	config := c.ParseConfig()

	Lock := new(sync.Mutex)
	Lock.Lock()
	chkpoint.SetChkpt(config.Chkpoint,logfile,"Offset",str_offset)
	Lock.Unlock()
	return string(data)
}


func ReadData(chkpt,logfile string)string{
	offset := chkpoint.Chkpt(chkpt,logfile,"Offset")
	fileinfo,err := os.Stat(logfile)
	if err != nil{log.Fatal(err)}
	size := fileinfo.Size()

	if offset <= size{
		data := TailLog(offset,logfile)
		return data
	}else{
		data := TailLog(0,logfile)
		return data
	}
}

//保存文件列表
var filelist []string

func multifiles(path string) []string{
	filepath.Walk(path,func(path string,fileinfo os.FileInfo,err error) error {
		if fileinfo == nil {
			return err
		}
		if fileinfo.IsDir() {
			return nil
		}

		filelist = append(filelist,path)
		return nil
	})
	return filelist
}


type Result struct{
	Path string
	Data string
}

//从文件头开始读
func ReadBigFile(file string,fd *os.File,size int64) string{
	c := common.Config{}
	config := c.ParseConfig()

	for {
		bufset := chkpoint.GetChkpt(config.Chkpoint,file,"Bufset")

		fd.Seek(bufset,0)
		//每次读取Buffer字节的内容
		Buffer := config.Buffer
		s := make([]byte,Buffer)

		nr,err := fd.Read(s)
		if err != nil || err == io.EOF{break}
		//记录读取到大文件的位置
		Lock := new(sync.Mutex)
		Lock.Lock()
		chkpoint.SetChkpt(config.Chkpoint,file,"Bufset",strconv.FormatInt(int64(int(bufset) + nr),10))
		Lock.Unlock()


		return string(s)
	}
	//一个大文件从头读到当前位置后记录下offset,后面开始实时读取
	Lock := new(sync.Mutex)
	Lock.Lock()
	chkpoint.SetChkpt(config.Chkpoint,file,"Offset",strconv.FormatInt(size,10))
	Lock.Unlock()
	return ""
}

func SendLog() chan Result{
	c := common.Config{}
	config := c.ParseConfig()

	files := multifiles(config.Logpath)
	ch := make(chan Result,len(files))

	var wg sync.WaitGroup
	wg.Add(1)

	//多个日志文件并发处理
	go func() {
		defer wg.Done()
		for _,file := range files{
			offset := chkpoint.Chkpt(config.Chkpoint,file,"Offset")

			fileinfo,_ := os.Stat(file)
			filename := fileinfo.Name()
			size := fileinfo.Size()

			if offset == 0{
				//从文件头开始读，可处理大文件
				fd,err := os.Open(file)
				if err != nil{log.Fatal(err)}
				defer fd.Close()

				data := ReadBigFile(file,fd,size)
				ch <- Result{filename,data}
			} else {
				//从文件结尾处实时读取
				if strings.Contains(file,".swp"){
					continue
				}
				data := ReadData(config.Chkpoint,file)
				ch <- Result{filename, data}
			}
		}
	}()

	wg.Wait()
	close(ch)
	return ch
}