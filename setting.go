package main

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/zhnxin/golang-image/jpeg"
	"github.com/zhnxin/mobi"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const defaultCoverWidth = 860
const defaultCoverHight = 1200

type Config struct {
	Title           string
	Cover           string
	CoverWidth      uint `toml:"cover_width"`
	CoverHight      uint `toml:"cover_hight"`
	Author          string
	Chapter         string
	SubChapter      string
	Encoding        string
	File            string
	ChapterRegex    *regexp.Regexp
	SubChapterRegex *regexp.Regexp
	Compress        bool
	decode          *encoding.Decoder
	Lang            string
}

func NewConfig(title, cover, thumbnail, author, chapter, subchapter, encoding, file string, compress bool) *Config {
	config := &Config{
		Title:      title,
		Cover:      cover,
		CoverWidth: defaultCoverWidth, CoverHight: defaultCoverHight,
		Author:          author,
		Chapter:         chapter,
		SubChapter:      subchapter,
		Encoding:        encoding,
		File:            file,
		Compress:        compress,
		ChapterRegex:    nil,
		SubChapterRegex: nil,
		decode:          nil,
		Lang:            "",
	}
	return config
}
func (c *Config) Update(file, title, author, cover, thumbnail string) {
	if file != "" {
		c.File = file
	}
	if title != "" {
		c.Title = title
	}
	if author != "" {
		c.Author = author
	}
	if cover != "" {
		c.Cover = cover
	}
}

func (config *Config) Check() (err error) {
	switch config.Encoding {
	case "GB18030", "gb18030":
		config.decode = simplifiedchinese.GB18030.NewDecoder()
	case "GBK", "gbk":
		config.decode = simplifiedchinese.GBK.NewDecoder()
	case "UTF8", "utf8", "utf-8", "":
		config.decode = nil
	default:
		return fmt.Errorf("Unsupport encoding[GB18030,GBK,UTF8(default)]:%s", config.Encoding)
	}
	if _, err = os.Stat(config.File); os.IsNotExist(err) {
		return
	}
	config.ChapterRegex, err = regexp.Compile(config.Chapter)
	if err == nil && config.SubChapter != "" {
		config.SubChapterRegex, err = regexp.Compile(config.SubChapter)
	}
	if config.Cover == "" {
		return errors.New("ebook need a cover")
	}
	return
}

func NewConfigWithFile(configFile string) (config *Config, err error) {
	config = &Config{}
	_, err = toml.DecodeFile(configFile, &config)
	if err != nil {
		return
	}
	if config.CoverWidth == 0 && config.CoverHight == 0 {
		config.CoverHight = defaultCoverHight
		config.CoverWidth = defaultCoverWidth
	}
	if config.CoverWidth == 0 || config.CoverHight == 0 {
		return nil, errors.New("CoverWidth or CoverHight should not be zero")
	}
	return
}

func (config *Config) NewWriter(fileName string) (mobi.Builder, error) {
	m := mobi.NewBuilder()
	m.Title(config.Title)
	if !config.Compress {
		m.Compression(mobi.CompressionNone)
	}
	if config.Lang != "" {
		m.NewExthRecord(mobi.EXTH_LANGUAGE, config.Lang)
	}
	if config.Cover != "" {
		cover, thumb, err := config.generateCover(config.Cover)
		if err != nil {
			return nil, err
		}
		m.AddCover(cover, thumb)
		defer func() {
			_ = os.Remove(cover)
			_ = os.Remove(thumb)
		}()
	}
	m.NewExthRecord(mobi.EXTH_DOCTYPE, "EBOK")
	m.NewExthRecord(mobi.EXTH_AUTHOR, config.Author)
	return m, nil
}

func (c *Config) Decode(content []byte) ([]byte, error) {
	if c.decode != nil {
		return c.decode.Bytes(content)
	}
	return content, nil
}

func (c *Config) generateCover(imgPath string) (cover string, thumb string, err error) {
	var img image.Image
	img, err = imaging.Open(imgPath)
	if err != nil {
		return
	}
	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()
	ts := fmt.Sprintf("%d", time.Now().Unix())
	cover = ts + ".jpg"
	thumb = ts + "_thumb.jpg"
	if float32(dx)/float32(dy) >= float32(c.CoverWidth)/float32(c.CoverHight) {
		dst := resize.Resize(c.CoverWidth, 0, img, resize.Lanczos3)
		err = saveImg(dst, cover)
		if err != nil {
			return
		}
		dst = resize.Resize(uint(float32(c.CoverWidth)*0.3), 0, img, resize.Lanczos3)
		err = saveImg(dst, thumb)
		return
	}
	dst := resize.Resize(0, c.CoverHight, img, resize.Lanczos3)
	err = saveImg(dst, cover)
	if err != nil {
		return
	}
	dst = resize.Resize(0, uint(float32(c.CoverHight)*0.3), img, resize.Lanczos3)
	err = saveImg(dst, thumb)
	return
}

func saveImg(img image.Image, filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	return jpeg.EncodeWithJfif(out, img, nil, nil)
}
