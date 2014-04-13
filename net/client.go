package net

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"pardal/protocol"
)

const (
	secret = "sheacsrhet"
	uri    = "http://sinespcidadao.sinesp.gov.br/sinesp-cidadao/ConsultaPlaca"
)

type SinespClient struct {
	http.Client
}

func NewSinespClient() *SinespClient {
	c := new(SinespClient)
	c.Client = http.Client{}
	// Golang's http client uses HTTP 1.1 but sinesp server sends back a
	// "Connection: close", unfortunately. Therefore the socket cannot be
	// reused.
	return c
}

func (c *SinespClient) GetVehicleInfo(plate string) (*protocol.Vehicle, error) {
	token := fmt.Sprintf("%x",
		string(hmacSHA1([]byte(plate), []byte(secret))))

	resp, err := c.Post(uri, "text/plain", buildBody(plate, token))

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return protocol.Unmarshal(body)
}

func hmacSHA1(input, key []byte) []byte {
	h := hmac.New(sha1.New, key)
	h.Write(input)
	return h.Sum(nil)
}

func buildBody(placa, token string) io.Reader {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(`<v:Envelope xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns:d="http://www.w3.org/2001/XMLSchema" xmlns:c="http://schemas.xmlsoap.org/soap/encoding/" xmlns:v="http://schemas.xmlsoap.org/soap/envelope/"><v:Header><dispositivo>motorola MB525</dispositivo><nomeSO>Android1.1.1</nomeSO><versaoSO>2.2.1</versaoSO><aplicativo>aplicativo</aplicativo><ip>192.168.10.43</ip><token>`)
	buf.WriteString(token)
	buf.WriteString(`</token><latitude>0.1</latitude><longitude>0.0</longitude><versaoAplicativo /></v:Header><v:Body><n0:getStatus xmlns:n0="http://soap.ws.placa.service.sinesp.serpro.gov.br/"><placa>`)
	buf.WriteString(placa)
	buf.WriteString(`</placa></n0:getStatus></v:Body></v:Envelope>`)

	return buf
}
