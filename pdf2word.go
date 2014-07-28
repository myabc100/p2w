package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"
)

const (
	exe  = `F:\godir\src\pdf2word\pdf2word-v3.0\pdf2word.exe`
	pdf  = `F:\godir\src\pdf2word\pdf.pdf`
	word = `F:\godir\src\pdf2word\pdf.doc`
)

func main() {
	http.HandleFunc("/api/upload.aspx", upload)
	err := http.ListenAndServe(":81", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	swith()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("tpl/upload.tpl")
		t.Execute(w, nil)
	} else {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Add("content-type", "text/html")
				io.WriteString(w, err.(string))
			}
		}()

		er := r.ParseMultipartForm(1000000)

		if er != nil {
			w.Header().Add("content-type", "text/html")
			io.WriteString(w, er.Error())
			return
		}

	
		for _, v := range r.MultipartForm.File {
			for _, uf := range v {

				filename := uf.Filename
				//结束文件
				file, err := uf.Open()
				if err != nil {
					fmt.Println(err)
				}
				//保存文件
				defer file.Close()
				

				f, err := os.Create(getpath(filename))
				defer f.Close()
				io.Copy(f, file)
				fstat, _ := f.Stat()

				fmt.Fprintf(w, " NO.: %d  Size: %d KB  Name：%s\n", time.Now().Format("2006-01-02 15:04:05"), fstat.Size()/1024, filename)
				
			}

		}

	}
}

func getpath(fp string) nfn string {
	t := time.Now()
	fpath := fmt.Sprintf("updata/%s/%s/%s/%s%s", t.Month(), t.Day(), t.Hour())
	os.MkdirAll(fpath, "0777")
	newfn = fmt.Sprintf("%s/", fpath, t.Unix(), path.Ext(fp))
}

/*
func down(w http.ResponseWriter, r *http.Request) {

}*/

func to(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//fn := r.Form["fn"]
	c := make(chan int, 1)

	cmd := exec.Command(exe, "-q", "-m", "-r", "-i", pdf, "-o", word)
	cmd.Start()
	go func() {
		cmd.Wait()
		c <- 1
	}()

	select {
	case <-time.After(3 * time.Second):
		cmd.Process.Kill()
		fmt.Println("timeout")
	case <-c:
		fmt.Println("done")
	}

	t := <-time.After(3 * time.Second)
	fmt.Println(t)
}

/*
func tocheck(w http.ResponseWriter, r *http.Request) {

}*/
