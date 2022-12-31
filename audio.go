package cob

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Audio struct {

}

var DAudio = &Audio{}

// ConvertMggToMp3 格式化文件名、解密、转换一步到位
func (obj *Audio)ConvertMggToMp3(umPath string, ffmpegPath string, mggPath string, oggPath string, mp3Path string)  {
	obj.DeleteFilenameBlank(mggPath)
	obj.ConvertMggToOgg(umPath, mggPath, oggPath)
	obj.ConvertOggToMp3(ffmpegPath, oggPath, mp3Path)
}

// DeleteFilenameBlank 遍历目录，删除文件名中的空格并重命名
func (obj *Audio)DeleteFilenameBlank(root string)  {
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			newPath := filepath.Join(root, strings.ReplaceAll(d.Name(), " ", ""))
			//fmt.Println(path)
			//fmt.Println(newPath)
			err = os.Rename(path, newPath)
			if err != nil {
				log.Fatalln(err)
			}
		}
		return nil
	})
}

// ConvertOggToMp3 通过ffmpeg将ogg转换为mp3格式
func (obj *Audio)ConvertOggToMp3(ffmpegPath string, oggPath string, mp3Path string)  {
	filepath.WalkDir(oggPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			newPath := filepath.Join(mp3Path, d.Name()[:len(d.Name())-3]+"mp3")
			cmd := exec.Command(ffmpegPath, "-i", path, "-acodec", "libmp3lame", newPath)
			//cmd := exec.Command(ffmpegPath, "-i", path, "-acodec", "mp3", newPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatalln(err)
			}
		}
		return nil
	})
}

// ConvertMggToOgg 通过um将mgg格式解密为ogg格式
func (obj *Audio)ConvertMggToOgg(umPath string, mggPath string, oggPath string)  {
	cmd := exec.Command(umPath, "-o", oggPath, "-i", mggPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}
