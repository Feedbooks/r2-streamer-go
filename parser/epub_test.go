package parser

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type testDataStruct struct {
	filepath     string
	err          error
	title        string
	authorName   string
	identifier   string
	thirdChapter string
	tocChildren  bool
	Source       string
	NoLinear     string
	MultipleLang bool
}

func TestPublication(t *testing.T) {
	testData := []testDataStruct{
		{"../test/empty.epub", errors.New("can't open or parse epub file with err : open ../test/empty.epub: no such file or directory"), "", "", "", "", false, "", "", false},
		{"../test/moby-dick.epub", nil, "Moby-Dick", "Herman Melville", "code.google.com.epub-samples.moby-dick-basic", "ETYMOLOGY.", false, "", "cover.xhtml", false},
		{"../test/kusamakura.epub", nil, "草枕", "夏目 漱石", "http://www.aozora.gr.jp/cards/000148/card776.html", "三", false, "", "", true},
		{"../test/feedbooks_book_6816.epub", nil, "Mémoires d'Outre-tombe", "François-René de Chateaubriand", "urn:uuid:47f6aaf6-aa7e-11e6-8357-4c72b9252ec6", "Partie 1", true, "www.ebooksfrance.com", "", false},
	}

	for _, d := range testData {
		Convey("Given "+d.title+" book", t, func() {
			publication, err := Parse(d.filepath)
			Convey("There no exception parsing", func() {
				if d.err != nil {
					So(err.Error(), ShouldEqual, d.err.Error())
				} else {
					So(err, ShouldEqual, nil)
				}
			})

			if d.MultipleLang == true {
				Convey("The title has multiple language", func() {
					So(publication.Metadata.Title.MultiString, ShouldNotBeEmpty)
				})

				fmt.Println(publication.Metadata.Language[0])
				Convey("The title is good", func() {
					So(publication.Metadata.Title.MultiString[publication.Metadata.Language[0]], ShouldEqual, d.title)
				})
			} else {
				Convey("The title is good", func() {
					So(publication.Metadata.Title.String(), ShouldEqual, d.title)
				})
			}

			if err == nil {
				Convey("There must be an author", func() {
					So(len(publication.Metadata.Author), ShouldBeGreaterThanOrEqualTo, 1)
				})
			}

			if d.authorName != "" && len(publication.Metadata.Author) > 0 {
				Convey("first author is good", func() {
					So(publication.Metadata.Author[0].Name.String(), ShouldEqual, d.authorName)
				})
			}

			Convey("Identifier is good", func() {
				So(publication.Metadata.Identifier, ShouldEqual, d.identifier)
			})

			Convey("The third chapter is good", func() {
				if len(publication.TOC) > 3 {
					So(publication.TOC[2].Title, ShouldEqual, d.thirdChapter)
				}
			})

			Convey("There Chapter with children", func() {
				emptyChildren := false
				for _, toc := range publication.TOC {
					if len(toc.Children) > 0 {
						emptyChildren = true
					}
				}
				if d.tocChildren == true {
					So(emptyChildren, ShouldBeTrue)
				} else {
					So(emptyChildren, ShouldBeFalse)
				}
			})

			Convey("dc:source is good", func() {
				So(publication.Metadata.Source, ShouldEqual, d.Source)
			})

			if d.NoLinear != "" {
				Convey("item no linear is not in spine", func() {
					findItemInSpine := false

					for _, it := range publication.Spine {
						if it.Href == d.NoLinear {
							findItemInSpine = true
						}
					}

					So(findItemInSpine, ShouldEqual, false)
				})
			}

		})
	}

}
