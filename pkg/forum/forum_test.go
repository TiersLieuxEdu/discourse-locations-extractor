package forum

import (
	"math"
	"testing"

	"github.com/TiersLieuxEdu/discourse-locations-extractor/pkg/lieux"
)

func TestExtractInfoOnEmptySource(t *testing.T) {
	var info lieux.Info
	info.Name = ""
	ExtractInfo("", &info)
	if info.Name != "" {
		t.Errorf("Name shall not have been set. Got %s", info.Name)
	}
}

func TestExtractInfoSite(t *testing.T) {
	var info lieux.Info
	info.WebSite = ""
	ExtractInfo(`## Informations

<dl class="h-geo" id="info">
<dt>Site</dt><dd>https://www.example.com</dd>
<dt>adresse</dt><dd>42 rue des communs<br/>100000  TiersLieuxEdu</dd>
<dt>Latitude</dt><dd>43.130</dd>
<dt>Longitude</dt><dd>5.936</dd>
<dt>Tags</dt><dd>#tag1, #tag2</dd>
</dl>`, &info)
	if info.WebSite != "https://www.example.com" {
		t.Errorf("Name shall have been set to https://www.example.com. Got %s", info.WebSite)
	}
}

func TestExtractInfoFloatLatitude(t *testing.T) {
	var info lieux.Info
	info.Latitude = ""
	info.Lat = 10.0
	ExtractInfo(`## Informations

<dl class="h-geo" id="info">
<dt>Site</dt><dd>https://www.example.com</dd>
<dt>adresse</dt><dd>42 rue des communs<br/>100000  TiersLieuxEdu</dd>
<dt>Latitude</dt><dd>43.130</dd>
<dt>Longitude</dt><dd>5.936</dd>
<dt>Tags</dt><dd>#tag1, #tag2</dd>
</dl>`, &info)
	if info.Lat != 43.130 {
		t.Errorf("Latitude shall have been set to 43.130. Got %f (%s)", info.Lat, info.Latitude)
	}
}

func TestExtractInfoFloatLongitude(t *testing.T) {
	var info lieux.Info
	info.Longitude = "0.0"
	info.Long = 10.0
	ExtractInfo(`## Informations

<dl class="h-geo" id="info">
<dt>Site</dt><dd>https://www.example.com</dd>
<dt>adresse</dt><dd>42 rue des communs<br/>100000  TiersLieuxEdu</dd>
<dt>Latitude</dt><dd>43.130</dd>
<dt>Longitude</dt><dd>5.936</dd>
<dt>Tags</dt><dd>#tag1, #tag2</dd>
</dl>`, &info)
	if info.Long != 5.936 {
		t.Errorf("Longitude shall have been set to 5.936. Got %f (%s)", info.Long, info.Longitude)
	}
}

func TestExtractInfoDummyLatitude(t *testing.T) {
	var info lieux.Info
	info.Latitude = "0.0"
	info.Lat = 10.0
	ExtractInfo(`## Informations

<dl class="h-geo" id="info">
<dt>Site</dt><dd>https://www.example.com</dd>
<dt>adresse</dt><dd>42 rue des communs<br/>100000  TiersLieuxEdu</dd>
<dt>Latitude</dt><dd>TBD</dd>
<dt>Longitude</dt><dd>TBD</dd>
<dt>Tags</dt><dd>#tag1, #tag2</dd>
</dl>`, &info)
	if info.Lat != 0.0 {
		t.Errorf("Latitude shall be reset to zero. Got %f (%s)", info.Lat, info.Latitude)
	}
}

func TestExtractInfoDummyLongitude(t *testing.T) {
	var info lieux.Info
	info.Longitude = "0.0"
	info.Long = 10.0
	ExtractInfo(`## Informations

<dl class="h-geo" id="info">
<dt>Site</dt><dd>https://www.example.com</dd>
<dt>adresse</dt><dd>42 rue des communs<br/>100000  TiersLieuxEdu</dd>
<dt>Latitude</dt><dd>TBD</dd>
<dt>Longitude</dt><dd>TBD</dd>
<dt>Tags</dt><dd>#tag1, #tag2</dd>
</dl>`, &info)
	if info.Long != 0.0 {
		t.Errorf("Longitude shall be reset to zero. Got %f (%s)", info.Long, info.Longitude)
	}
}

func TestExtractInfoDegreesMinutesSecondsLatitude(t *testing.T) {
	var info lieux.Info
	info.Latitude = "0.0"
	info.Lat = 10.0
	ExtractInfo(`## Informations

<dl class="h-geo" id="info">
<dt>Site</dt><dd>https://www.example.com</dd>
<dt>adresse</dt><dd>42 rue des communs<br/>100000  TiersLieuxEdu</dd>
<dt>Latitude</dt><dd>50°44'41.9"N</dd>
<dt>Longitude</dt><dd>3°13'15.8"E</dd>
<dt>Tags</dt><dd>#tag1, #tag2</dd>
</dl>`, &info)
	if math.Abs(info.Lat-50.744972) > 0.000001 {
		t.Errorf("Latitude shall be converted to 50.744972. Got %f (%s)", info.Lat, info.Latitude)
	}
}

func TestExtractInfoDegreesMinutesSecondsLongitude(t *testing.T) {
	var info lieux.Info
	info.Longitude = "0.0"
	info.Long = 10.0
	ExtractInfo(`## Informations

<dl class="h-geo" id="info">
<dt>Site</dt><dd>https://www.example.com</dd>
<dt>adresse</dt><dd>42 rue des communs<br/>100000  TiersLieuxEdu</dd>
<dt>Latitude</dt><dd>50°44'41.9"N</dd>
<dt>Longitude</dt><dd>3°13'15.8"E</dd>
<dt>Tags</dt><dd>#tag1, #tag2</dd>
</dl>`, &info)
	if math.Abs(info.Long-3.221056) > 0.000001 {
		t.Errorf("Longitude shall be converted to 3.221056. Got %f (%s)", info.Long, info.Longitude)
	}
}

func TestExtractInfoNegativeDegreesMinutesSecondsLatitude(t *testing.T) {
	var info lieux.Info
	info.Latitude = "0.0"
	info.Lat = 10.0
	ExtractInfo(`## Informations

<dl class="h-geo" id="info">
<dt>Site</dt><dd>https://www.example.com</dd>
<dt>adresse</dt><dd>42 rue des communs<br/>100000  TiersLieuxEdu</dd>
<dt>Latitude</dt><dd>50°44'41.9"S</dd>
<dt>Longitude</dt><dd>3°13'15.8"W</dd>
<dt>Tags</dt><dd>#tag1, #tag2</dd>
</dl>`, &info)
	if math.Abs(info.Lat+50.744972) > 0.000001 {
		t.Errorf("Latitude shall be converted to -50.744972. Got %f (%s)", info.Lat, info.Latitude)
	}
}

func TestExtractInfoNegativeDegreesMinutesSecondsLongitude(t *testing.T) {
	var info lieux.Info
	info.Longitude = "0.0"
	info.Long = 10.0
	ExtractInfo(`## Informations

<dl class="h-geo" id="info">
<dt>Site</dt><dd>https://www.example.com</dd>
<dt>adresse</dt><dd>42 rue des communs<br/>100000  TiersLieuxEdu</dd>
<dt>Latitude</dt><dd>50°44'41.9"S</dd>
<dt>Longitude</dt><dd>3°13'15.8"W</dd>
<dt>Tags</dt><dd>#tag1, #tag2</dd>
</dl>`, &info)
	if math.Abs(info.Long+3.221056) > 0.000001 {
		t.Errorf("Longitude shall be converted to -3.221056. Got %f (%s)", info.Long, info.Longitude)
	}
}

func TestConvertPositionning(t *testing.T) {
	_, err := ConvertPositionning(` 50°44'41.9"N`)
	if err != nil {
		t.Errorf("Position could not be converted %s", err)
	}
}

func TestExtract3Tags(t *testing.T) {
	var info lieux.Info
	info.Longitude = "0.0"
	info.Long = 10.0
	info.Tags = make([]string, 0, 6)
	ExtractInfo(`## Informations

<dl class="h-geo" id="info">
<dt>Site</dt><dd>https://www.example.com</dd>
<dt>adresse</dt><dd>42 rue des communs<br/>100000  TiersLieuxEdu</dd>
<dt>Latitude</dt><dd>50°44'41.9"N</dd>
<dt>Longitude</dt><dd>3°13'15.8"E</dd>
<dt>Tags</dt><dd>#tag1, #tag2, tag3</dd>
</dl>`, &info)
	if len(info.Tags) != 3 {
		t.Errorf("Shall found 3 tags, got only %d: %s", len(info.Tags), info.Tags)
	}
}
