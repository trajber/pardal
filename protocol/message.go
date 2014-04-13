package protocol

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
)

var (
	ErrVehicleNotFound = errors.New("Vehicle not found")
)

type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Soap    *SoapBody
}

type SoapBody struct {
	XMLName           xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	GetStatusResponse *GetStatusResponse
}

type GetStatusResponse struct {
	XMLName xml.Name `xml:"http://soap.ws.placa.service.sinesp.serpro.gov.br/ getStatusResponse"`
	Vehicle *Vehicle `xml:"return"`
}

type Vehicle struct {
	CodigoRetorno   int    `xml:"codigoRetorno"`
	MensagemRetorno string `xml:"mensagemRetorno"`
	Placa           string `xml:"placa"`
	CodigoSituacao  int    `xml:"codigoSituacao"`
	Situacao        string `xml:"situacao"`
	Modelo          string `xml:"modelo"`
	Marca           string `xml:"marca"`
	Cor             string `xml:"cor"`
	Ano             int    `xml:"ano"`
	AnoModelo       int    `xml:"anoModelo"`
	Data            string `xml:"data"`
	UF              string `xml:"uf"`
	Municipio       string `xml:"municipio"`
	Chassi          string `xml:"chassi"`
}

func (v *Vehicle) String() string {
	buf := newPrinter()
	buf.WriteLabelAndMessage("Placa", v.Placa)
	buf.WriteLabelAndMessage("Modelo", v.Modelo)
	buf.WriteLabelAndMessage("Marca", v.Marca)
	buf.WriteLabelAndMessage("Cor", v.Cor)
	buf.WriteLabelAndMessage("Ano", v.Ano)
	buf.WriteLabelAndMessage("Ano/Modelo", v.AnoModelo)
	buf.WriteLabelAndMessage("UF", v.UF)
	buf.WriteLabelAndMessage("Município", v.Municipio)
	buf.WriteLabelAndMessage("Chassi", v.Chassi)
	return buf.String()
}

func Unmarshal(data []byte) (*Vehicle, error) {
	// the response is ISO-8859-1, we need to convert it to UTF-8
	data = toUTF8(data)

	var v Envelope
	if err := xml.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	// TODO handle all possible errors
	if v.Soap.GetStatusResponse.Vehicle.CodigoRetorno == 3 {
		return nil, ErrVehicleNotFound
	}

	return v.Soap.GetStatusResponse.Vehicle, nil
}

func toUTF8(input []byte) []byte {
	buf := bytes.NewBuffer(nil)
	for _, b := range input {
		buf.WriteRune(rune(b))
	}
	return buf.Bytes()
}

func newPrinter() *printer {
	p := new(printer)
	p.Buffer = bytes.NewBuffer(nil)
	return p
}

type printer struct {
	*bytes.Buffer
}

func (p *printer) WriteLabelAndMessage(label string, msg interface{}) {
	p.WriteString(label)
	p.WriteString(":")
	switch msg.(type) {
	case string:
		p.WriteString(msg.(string))
	case int:
		p.WriteString(fmt.Sprintf("%d", msg))
	}
	p.WriteByte('\n')
}
