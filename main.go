package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/zhnxin/mobi"
)

var (
	HELP           = flag.Bool("h", false, "help")
	ConfigFile     = flag.String("config", "", "ebook config file(.toml)")
	IsParagraph    = flag.Bool("p", false, "[option]is to use <p></p>,use false as default")
	OutputFileName = flag.String("o", "", "[option]output file name")
	IsEscape       = flag.Bool("escape", false, "[option]To Disable html escape")
	MetaFile       = flag.String("f", "", "input file")
	MetaCover      = flag.String("cover", "", "mobi cover")
	MetaThumb      = flag.String("thumb", "", "mobi thumbnail")
	MetaTitle      = flag.String("title", "", "mobi title")
	MetaAuthor     = flag.String("author", "", "EBOK author")
	MetaCompress   = flag.Bool("compress", false, "Is to compress")
	MetaEncoding   = flag.String("encoding", "gb18030", "encoding:gb18030(default),gbk,uft-8")
	MetaChapter    = flag.String("chapter", "^第[零一二三四五六七八九十百千两\\d]+章[　 ]{0,1}.*$", "regexp pattern for chapter,default:'^第[零一二三四五六七八九十百千两\\d]+章 .*$'")
	MetaSubChapter = flag.String("subchapter", "", "regexp pattern for chapter,default:'^第[零一二三四五六七八九十百千两\\d]+章[　 ]{0,1}.*$'")

	CONFIG *Config
)

func initMain() {
	flag.Parse()
	var err error
	if *ConfigFile != "" {
		CONFIG, err = NewConfigWithFile(*ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		CONFIG.Update(*MetaFile, *MetaTitle, *MetaAuthor, *MetaCover, *MetaThumb)
	} else {
		CONFIG = NewConfig(*MetaTitle, *MetaCover, *MetaCover, *MetaAuthor, *MetaChapter, *MetaSubChapter, *MetaEncoding, *MetaFile, *MetaCompress)
	}

	if err = CONFIG.Check(); err != nil {
		log.Fatal(err)
		flag.Usage()
		return
	}
}

func main() {
	initMain()
	if *HELP {
		flag.Usage()
		fmt.Println(`Sugesstion:
	chapter: "^第[零一二三四五六七八九十百千两\\d]+[卷部][　 ]{0,1}.*$"
	subchapter: "^第[零一二三四五六七八九十百千两\\d]+章[　 ]{0,1}.*$"`)
		return
	}
	SetIsEscape(*IsEscape)
	SetIsParagraph(*IsParagraph)
	mobi.SetSkipLog(true)

	file, err := os.Open(CONFIG.File)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	mobiWriter, err := CONFIG.NewWriter(*OutputFileName)
	if err != nil {
		log.Fatal(err)
	}

	chapter := NewChapter(CONFIG.Title)
	var line []byte
	for scanner.Scan() {
		line, err = CONFIG.Decode(scanner.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		if CONFIG.ChapterRegex.Match(line) {
			chapter.Flush(mobiWriter)
			chapter.Restore(string(line))
		} else if CONFIG.SubChapterRegex != nil && CONFIG.SubChapterRegex.Match(line) {
			chapter.AddSubChapter(string(line))
		} else {
			chapter.Append(line)
		}
	}
	chapter.Flush(mobiWriter)
	outfile, err := os.OpenFile(CONFIG.Title+".mobi", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	mobiWriter.WriteTo(outfile)
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
