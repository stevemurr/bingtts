package bingtts

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

var voices = map[string]string{
	"ar-eg female": "Microsoft Server Speech Text to Speech Voice (ar-EG, Hoda)",
	"de-de female": "Microsoft Server Speech Text to Speech Voice (de-DE, Hedda)",
	"de-de male":   "Microsoft Server Speech Text to Speech Voice (de-DE, Stefan, Apollo)",
	"en-au female": "Microsoft Server Speech Text to Speech Voice (en-AU, Catherine)",
	"en-ca female": "Microsoft Server Speech Text to Speech Voice (en-CA, Linda)",
	"en-gb female": "Microsoft Server Speech Text to Speech Voice (en-GB, Susan, Apollo)",
	"en-gb male":   "Microsoft Server Speech Text to Speech Voice (en-GB, George, Apollo)",
	"en-in male":   "Microsoft Server Speech Text to Speech Voice (en-IN, Ravi, Apollo)",
	"en-us female": "Microsoft Server Speech Text to Speech Voice (en-US, ZiraRUS)",
	"en-us male":   "Microsoft Server Speech Text to Speech Voice (en-US, BenjaminRUS)",
	"es-es female": "Microsoft Server Speech Text to Speech Voice (es-ES, Laura, Apollo)",
	"es-es male":   "Microsoft Server Speech Text to Speech Voice (es-ES, Pablo, Apollo)",
	"es-mx male":   "Microsoft Server Speech Text to Speech Voice (es-MX, Raul, Apollo)",
	"fr-ca female": "Microsoft Server Speech Text to Speech Voice (fr-CA, Caroline)",
	"fr-fr female": "Microsoft Server Speech Text to Speech Voice (fr-FR, Julie, Apollo)",
	"fr-fr male":   "Microsoft Server Speech Text to Speech Voice (fr-FR, Paul, Apollo)",
	"it-it male":   "Microsoft Server Speech Text to Speech Voice (it-IT, Cosimo, Apollo)",
	"ja-jp female": "Microsoft Server Speech Text to Speech Voice (ja-JP, Ayumi, Apollo)",
	"ja-jp male":   "Microsoft Server Speech Text to Speech Voice (ja-JP, Ichiro, Apollo)",
	"pt-br male":   "Microsoft Server Speech Text to Speech Voice (pt-BR, Daniel, Apollo)",
	"ru-ru female": "Microsoft Server Speech Text to Speech Voice (ru-RU, Irina, Apollo)",
	"ru-ru male":   "Microsoft Server Speech Text to Speech Voice (ru-RU, Pavel, Apollo)",
	"zh-cn female": "Microsoft Server Speech Text to Speech Voice (zh-CN, Yaoyao, Apollo)",
	"zh-cn male":   "Microsoft Server Speech Text to Speech Voice (zh-CN, Kangkang, Apollo)",
	"zh-hk female": "Microsoft Server Speech Text to Speech Voice (zh-HK, Tracy, Apollo)",
	"zh-hk male":   "Microsoft Server Speech Text to Speech Voice (zh-HK, Danny, Apollo)",
	"zh-tw female": "Microsoft Server Speech Text to Speech Voice (zh-TW, Yating, Apollo)",
	"zh-tw male":   "Microsoft Server Speech Text to Speech Voice (zh-TW, Zhiwei, Apollo)",
}

const (
	bingSpeechTokenEndpoint = "https://api.cognitive.microsoft.com/sts/v1.0/issueToken"
	bingSpeechEndpointTTS   = "https://speech.platform.bing.com/synthesize"
	// RIFF16Bit16kHzMonoPCM --
	RIFF16Bit16kHzMonoPCM = "riff-16khz-16bit-mono-pcm"
	// RIFF8Bit8kHzMonoPCM --
	RIFF8Bit8kHzMonoPCM = "riff-8khz-8bit-mono-mulaw"
	// RAW8Bit8kHzMonoMulaw --
	RAW8Bit8kHzMonoMulaw = "raw-8khz-8bit-mono-mulaw"
	// RAW16Bit16kHzMonoMulaw --
	RAW16Bit16kHzMonoMulaw = "raw-16khz-16bit-mono-pcm"
)

func getSSML(locale, font, gender, text string) string {
	return fmt.Sprintf(`<speak version='1.0' xml:lang='%s'><voice name='%s' xml:lang='%s' xml:gender='%s'>%s</voice></speak>`,
		locale,
		font,
		locale,
		gender,
		text)
}

// Synthesize --
func Synthesize(token, text, locale, gender, outputFormat string) ([]byte, error) {
	client := &http.Client{}
	font := voices[fmt.Sprintf("%s %s", locale, gender)]
	ssml := getSSML(locale, font, gender, text)
	req, err := http.NewRequest("POST", bingSpeechEndpointTTS, bytes.NewBufferString(ssml))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Length", strconv.Itoa(len(token)))
	req.Header.Add("Content-Type", "application/ssml+xml")
	req.Header.Add("X-Microsoft-OutputFormat", outputFormat)
	req.Header.Add("X-Search-AppId", "00000000000000000000000000000000")
	req.Header.Add("X-Search-ClientID", "00000000000000000000000000000000")
	req.Header.Add("User-Agent", "go-bing-tts")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", res.Status)
	}
	defer res.Body.Close()
	size, err := strconv.Atoi(res.Header.Get("Content-Length"))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0, size))
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetVoices -- Return voices availble on cognitive services
func GetVoices() map[string]string {
	return voices
}

// IssueToken -- Get a JWT token from cognitive services
func IssueToken(key string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", bingSpeechTokenEndpoint, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Ocp-Apim-Subscription-Key", key)
	req.Header.Add("Content-Length", "0")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s", res.Status)
	}

	defer res.Body.Close()
	size, err := strconv.Atoi(res.Header.Get("Content-Length"))
	if err != nil {
		return "", err
	}
	buf := make([]byte, size)
	res.Body.Read(buf)
	return string(buf), nil
}
