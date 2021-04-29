package v1Branch

import "encoding/xml"

type Annotation struct {
  XMLName   xml.Name   `xml:"annotation"`
  Folder    string     `xml:"folder"`
  FileName  string     `xml:"filename"`
  Path      string     `xml:"path,omitempty"`
  Source    *Source    `xml:"source"`
  Size      *Size      `xml:"size"`
  Segmented int        `xml:"segmented"`
  Object    [] *Object `xml:"object"`
}

type Source struct {
  DataBase string `xml:"database"`
}

type Size struct {
  Width  int `xml:"width"`
  Height int `xml:"height"`
  Depth  int `xml:"depth"`
}

type Object struct {
  Name      string  `xml:"name"`
  Pose      string  `xml:"pose"`
  Truncated int     `xml:"truncated"`
  Difficult int     `xml:"difficult"`
  BndBox    *BndBox `xml:"bndbox"`
}

type BndBox struct {
  XMin string `xml:"xmin"`
  YMin string `xml:"ymin"`
  XMax string `xml:"xmax"`
  YMax string `xml:"ymax"`
}
