package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func upload(w http.ResponseWriter, r *http.Request) {
	if res, err := upload2(r); err != nil {
		_, _ = w.Write([]byte(err.Error() + fmt.Sprintf("\n已上传%d字节", res)))
	} else {
		_, fileHeader, _ := r.FormFile("file")
		cmd := exec.Command(ffmpegPath, "-y", "-i", fileHeader.Filename, "-t", "10", "-vf", "scale='max(180,iw*0.1)':-1", "-r", "10", "result.gif")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(w, fmt.Sprint(err)+": "+stderr.String())
			return
		}
		w.Write([]byte(uploadWeb))
	}

}

func upload2(r *http.Request) (int64, error) {
	upfile, fileHeader, err := r.FormFile("file")
	if !supportType(fileHeader.Filename) {
		return 0, errors.New("请上传mp4/flv/avi视频文件")
	}
	if err != nil {
		return 0, errors.New("上传文件错误")
	}
	log.Printf("文件名称：%s\n文件大小：%d字节\n", fileHeader.Filename, fileHeader.Size)
	filePath := "./" + fileHeader.Filename
	fileBool, err := isFileExists(filePath)
	if fileBool && err == nil {
		fmt.Println("文件已存在")
	} else {
		newfile, err := os.Create(filePath)
		defer newfile.Close()
		if err != nil {
			return 0, errors.New("创建文件失败")
		}
	}
	fi, _ := os.Stat(filePath)
	if fi.Size() == fileHeader.Size {
		return fileHeader.Size, nil
	}
	start := strconv.Itoa(int(fi.Size()))
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		return 0, errors.New("打开之前上传文件不存在")
	}
	count, _ := strconv.ParseInt(start, 10, 64)
	if count == 0 {
		fmt.Println("是服务器未接收过的新文件")
	} else {
		fmt.Println("服务器已接收字节数:", count)
	}
	upfile.Seek(count, 0)
	file.Seek(count, 0)
	data := make([]byte, 1024, 1024)
	var upTotal int64 = 0
	for {
		total, err := upfile.Read(data)
		if err == io.EOF {
			fmt.Println("文件复制完毕")
			break
		}
		len, err := file.Write(data[:total])
		if err != nil {
			return count, errors.New("文件上传失败")
		}
		upTotal += int64(len)
		count += int64(len)
		if count > 1000000 && fileBool == false {
			return count, errors.New("模拟上传中断")
		}
	}
	fmt.Println("文件上传成功，字节数:", upTotal)
	return upTotal, nil
}

// 判断文件或文件夹是否存在
func isFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

func supportType(fileName string) bool {
	supportType := []string{".mp4", ".flv", ".avi"}
	for _, i := range supportType {
		if fileName[len(fileName)-4:] == i {
			return true
		}
	}
	return false
}

const uploadWeb = `<html>
<head>
<title>视频转GIF</title>
</head>
<body>
<img src="result.gif" alt="Result">
</body>
</html>`
