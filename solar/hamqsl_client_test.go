package solar

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHamQSLClientParse(t *testing.T) {
	var solar Solar
	err := xml.Unmarshal([]byte(SAMPLE_RESPONSE), &solar)
	assert.NoError(t, err)

	// Spot check various solar values
	assert.Equal(t, 70, solar.SolarData.SolarFlux)
	assert.Equal(t, 4, solar.SolarData.AIndex)
	assert.Equal(t, 1, solar.SolarData.KIndex)
	assert.Equal(t, "No Report", solar.SolarData.KIndexNt)
	assert.Equal(t, "A3.6", solar.SolarData.XRay)
	assert.Equal(t, 0, solar.SolarData.Sunspots)
	assert.Equal(t, float32(9.86), solar.SolarData.MUF)

	// Assert HF band conditions parsed correctly
	assert.Equal(t, "80m-40m", solar.SolarData.CalculatedConditions.Bands[0].Name)
	assert.Equal(t, "day", solar.SolarData.CalculatedConditions.Bands[0].Time)
	assert.Equal(t, "Fair", solar.SolarData.CalculatedConditions.Bands[0].Condition)

	// Assert VHF phenomenon parsed correctly
	assert.Equal(t, "E-Skip", solar.SolarData.CalculatedConditionsVHF.Phenomenon[2].Name)
	assert.Equal(t, "north_america", solar.SolarData.CalculatedConditionsVHF.Phenomenon[2].Location)
	assert.Equal(t, "Band Closed", solar.SolarData.CalculatedConditionsVHF.Phenomenon[2].Condition)
}

const SAMPLE_RESPONSE = `
<solar>
<solardata>
<source url="http://www.hamqsl.com/solar.html">N0NBH</source>
<updated> 07 May 2021 0818 GMT</updated>
<solarflux>70</solarflux>
<aindex> 4</aindex>
<kindex> 1</kindex>
<kindexnt>No Report</kindexnt>
<xray>A3.6</xray>
<sunspots>0</sunspots>
<heliumline> 95.4</heliumline>
<protonflux>19</protonflux>
<electonflux>851</electonflux>
<aurora> 1</aurora>
<normalization>1.99</normalization>
<latdegree>67.5</latdegree>
<solarwind>292.5</solarwind>
<magneticfield> 1.2</magneticfield>
<calculatedconditions>
<band name="80m-40m" time="day">Fair</band>
<band name="30m-20m" time="day">Fair</band>
<band name="17m-15m" time="day">Poor</band>
<band name="12m-10m" time="day">Poor</band>
<band name="80m-40m" time="night">Good</band>
<band name="30m-20m" time="night">Fair</band>
<band name="17m-15m" time="night">Poor</band>
<band name="12m-10m" time="night">Poor</band>
</calculatedconditions>
<calculatedvhfconditions>
<phenomenon name="vhf-aurora" location="northern_hemi">Band Closed</phenomenon>
<phenomenon name="E-Skip" location="europe">Band Closed</phenomenon>
<phenomenon name="E-Skip" location="north_america">Band Closed</phenomenon>
<phenomenon name="E-Skip" location="europe_6m">50MHz ES</phenomenon>
<phenomenon name="E-Skip" location="europe_4m">Band Closed</phenomenon>
</calculatedvhfconditions>
<geomagfield>VR QUIET</geomagfield>
<signalnoise>S0-S1</signalnoise>
<fof2>4.05</fof2>
<muffactor>2.43</muffactor>
<muf> 9.86</muf>
</solardata>
</solar>
`
