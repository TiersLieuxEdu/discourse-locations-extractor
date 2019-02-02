package lieux;

type Info struct {
  Name string // Name of the lab
  Latitude string // WGS84 X coordinate
  Lat float64
  Longitude string // WGS84 Y coordinate
  Long float64
  WebSite string // URL of the website of the lab
  Forum string // URL to the forum entry for the place
  //Tags []string
}

func (lieu Info) AsSlice() []string {
  r := make([]string, 5)
  r[0] = lieu.Name
  r[1] = lieu.Latitude
  r[2] = lieu.Longitude
  r[3] = lieu.WebSite
  r[4] = lieu.Forum
  return r
}
