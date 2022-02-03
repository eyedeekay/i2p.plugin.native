package shellservice

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/fuxingZhang/zip"
	"i2pgit.org/idk/reseed-tools/su3"
)

type PluginConfig struct {
	PluginName          *string   //0
	KeyName             *string   //1
	Signer              *string   //2
	Version             *string   //3
	License             *string   //4
	Date                *string   //5
	Author              *string   //6
	Website             *string   //7
	UpdateURL           *string   //8
	Description         *string   //9
	DescriptionLang     []*string //10
	ConsoleLinkName     *string   //11
	ConsoleLinkNameLang []*string //12
	ConsoleLinkURL      *string   //13
	ConsoleIcon         *string   //14
	ConsoleIconCode     *string   //15
	MinVersion          *string   //16
	MaxVerion           *string   //17
	MinJava             *string   //18
	MinJetty            *string   //19
	MaxJetty            *string   //20
	NoStop              *bool     //21
	NoStart             *bool     //22
	Restart             *bool     //23
	OnlyUpdate          *bool     //24
	OnlyInstall         *bool     //25
	ConsoleLinkTip      *string   //26
	ConsoleLinkTipLang  []*string //27
	SignerDirectory     *string   //28
}

func (pc *PluginConfig) Print() string {
	r := pc.PrintPluginName()  //0
	r += pc.PrintKeyName()     //1
	r += pc.PrintSigner()      //2
	r += pc.PrintVersion()     //3
	r += pc.PrintLicense()     //4
	r += pc.PrintDate()        //5
	r += pc.PrintAuthor()      //6
	r += pc.PrintWebsite()     //7
	r += pc.PrintUpdateURL()   //8
	r += pc.PrintDescription() //9
	//10
	r += pc.PrintConsoleLinkName() //11
	//12
	r += pc.PrintConsoleLinkURL()  //13
	r += pc.PrintConsoleIcon()     //14
	r += pc.PrintConsoleIconCode() //15
	r += pc.PrintMinVersion()      //16
	r += pc.PrintMaxVerion()       //17
	r += pc.PrintMinJava()         //18
	r += pc.PrintMinJetty()        //19
	r += pc.PrintMaxJetty()        //20
	r += pc.PrintNoStop()          //21
	r += pc.PrintNoStart()         //22
	r += pc.PrintRestart()         //23
	r += pc.PrintOnlyUpdate()      //24
	r += pc.PrintOnlyInstall()     //25
	r += pc.PrintConsoleLinkTip()  //26
	//27

	return r
}

func (pc *PluginConfig) PrintPluginName() string {
	if pc.PluginName == nil || *pc.PluginName == "" {
		log.Fatal("-name is a required field.")
	}
	return fmt.Sprintf("name=%s\n", *pc.PluginName)
}
func (pc *PluginConfig) PrintKeyName() string {
	if pc.KeyName == nil && *pc.KeyName == "" {
		return ""
	}
	return fmt.Sprintf("key=%s\n", *pc.KeyName)
}
func (pc *PluginConfig) PrintSigner() string {
	if pc.Signer == nil || *pc.Signer == "" {
		log.Fatal("-signer is a required field.")
	}
	return fmt.Sprintf("signer=%s\n", *pc.Signer)
}
func (pc *PluginConfig) PrintAuthor() string {
	if pc.Author == nil || *pc.Author == "" {
		if pc.Signer == nil || *pc.Signer == "" {
			log.Fatal("-signer is a required field.")
		}
		return fmt.Sprintf("author=%s\n", *pc.Signer)
	}
	return fmt.Sprintf("author=%s\n", *pc.Author)
}
func (pc *PluginConfig) PrintVersion() string {
	if pc.Version == nil || *pc.Version == "" {
		log.Fatal("-version is a required field.")
	}
	return fmt.Sprintf("version=%s\n", *pc.Version)
}
func (pc *PluginConfig) PrintDate() string {
	if pc.Date == nil || *pc.Date == "" {
		return ""
	}
	return fmt.Sprintf("date=%s\n", *pc.Date)
}

func (pc *PluginConfig) PrintNoStop() string {
	if pc.NoStop == nil {
		return fmt.Sprintf("disableStop=%t\n", false)
	}
	return fmt.Sprintf("disableStop=%t\n", *pc.NoStop)
}
func (pc *PluginConfig) PrintNoStart() string {
	if pc.NoStart == nil {
		return fmt.Sprintf("dont-start-at-install=%t\n", false)
	}
	return fmt.Sprintf("dont-start-at-install=%t\n", *pc.NoStart)
}
func (pc *PluginConfig) PrintRestart() string {
	if pc.Restart == nil {
		return fmt.Sprintf("router-restart-required=%t\n", false)
	}
	return fmt.Sprintf("router-restart-required=%t\n", *pc.Restart)
}
func (pc *PluginConfig) PrintOnlyUpdate() string {
	if pc.OnlyUpdate == nil {
		return fmt.Sprintf("update-only=%t\n", false)
	}
	return fmt.Sprintf("update-only=%t\n", *pc.OnlyUpdate)
}
func (pc *PluginConfig) PrintOnlyInstall() string {
	if pc.OnlyInstall == nil {
		return fmt.Sprintf("install-only=%t\n", false)
	}
	return fmt.Sprintf("install-only=%t\n", *pc.OnlyInstall)
}

func (pc *PluginConfig) PrintLicense() string {
	if pc.License == nil || *pc.License == "" {
		return fmt.Sprintf("license=%s\n", "unknown")
	}
	return fmt.Sprintf("license=%s\n", *pc.License)
}

func (pc *PluginConfig) PrintWebsite() string {
	if pc.Website == nil || *pc.Website == "" {
		return fmt.Sprintf("websiteURL=%s%s%s\n", "http://", *pc.PluginName, ".i2p")
	}
	return fmt.Sprintf("websiteURL=%s\n", *pc.Website)
}
func (pc *PluginConfig) PrintUpdateURL() string {
	if pc.UpdateURL == nil || *pc.UpdateURL == "" {
		return fmt.Sprintf("updateURL=%s%s%s%s%s\n", "http://", *pc.PluginName, ".i2p/", *pc.PluginName, ".su3")
	}
	if strings.HasSuffix(*pc.UpdateURL, "xpi2p") {
		return fmt.Sprintf("updateURL=%s\n", *pc.UpdateURL)
	}
	if strings.HasSuffix(*pc.UpdateURL, "su3") {
		return fmt.Sprintf("updateURL=%s\n", *pc.UpdateURL)
	}
	return fmt.Sprintf("updateURL=%s\n", *pc.UpdateURL)
}
func (pc *PluginConfig) PrintDescription() string {
	if pc.Description == nil || *pc.Description == "" {
		return fmt.Sprintf("description=\"%s\"\n", "Plugin config generated by i2p.plugin.native")
	}
	return fmt.Sprintf("description=%s\n", strings.Replace(strings.Replace(*pc.Description, "\n", "", -1), "\"", "", -1))
}

//func (pc *PluginConfig) PrintDescriptionLang() []string      { return []string{""} }
func (pc *PluginConfig) PrintConsoleLinkName() string {
	if pc.ConsoleLinkName == nil || *pc.ConsoleLinkName == "" {
		return ""
	}
	return fmt.Sprintf("consoleLinkName=%s\n", *pc.ConsoleLinkName)
}

func (pc *PluginConfig) PrintConsoleLinkTip() string {
	if pc.ConsoleLinkTip == nil || *pc.ConsoleLinkTip == "" {
		return ""
	}
	return fmt.Sprintf("consoleLinkTooltip=%s\n", *pc.ConsoleLinkTip)
}

func (pc *PluginConfig) PrintConsoleLinkURL() string {
	if pc.ConsoleLinkURL == nil || *pc.ConsoleLinkURL == "" {
		return ""
	}
	return fmt.Sprintf("consoleLinkURL=%s\n", *pc.ConsoleLinkURL)
}

//func (pc *PluginConfig) PrintConsoleLinkNameLang() []*string { return []string{""} }
func (pc *PluginConfig) PrintConsoleIcon() string {
	if pc.ConsoleIcon == nil || *pc.ConsoleIcon == "" {
		return ""
	}
	return fmt.Sprintf("console-icon=%s\n", *pc.ConsoleIcon)
}
func (pc *PluginConfig) PrintConsoleIconCode() string {
	if pc.ConsoleIconCode == nil || *pc.ConsoleIconCode == "" {
		return ""
	}
	bytes, err := ioutil.ReadFile(*pc.ConsoleIconCode)
	if err != nil {
		return ""
	}
	i2pbase64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-~")
	i2pb64 := i2pbase64.EncodeToString(bytes)
	return fmt.Sprintf("icon-code=%s\n", i2pb64)
}

func (pc *PluginConfig) PrintMinVersion() string {
	return ""
}
func (pc *PluginConfig) PrintMaxVerion() string {
	return ""
}
func (pc *PluginConfig) PrintMinJava() string {
	return ""
}
func (pc *PluginConfig) PrintMinJetty() string {
	return ""
}
func (pc *PluginConfig) PrintMaxJetty() string {
	return ""
}

func (pc *PluginConfig) CreateZip() error {
	err := zip.Dir("plugin", *pc.PluginName+".zip", false)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (pc *PluginConfig) CreateSu3() (*su3.File, error) {
	su3File := su3.New()
	su3File.FileType = su3.FileTypeZIP
	su3File.ContentType = su3.ContentTypePlugin
	su3File.Version = []byte(*pc.Version)

	err := pc.CreateZip()
	if err != nil {
		return nil, err
	}
	zipped, err := ioutil.ReadFile(*pc.PluginName + ".zip")
	if err != nil {
		return nil, err
	}
	su3File.Content = zipped

	su3File.SignerID = []byte(*pc.Signer)
	sk, err := pc.LoadPrivateKey(*pc.Signer)
	if err != nil {
		return nil, err
	}
	su3File.Sign(sk)
	return su3File, nil
}

func (pc *PluginConfig) LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keys, err := pc.keysPath(path)
	if err != nil {
		return nil, err
	}
	privPem, err := ioutil.ReadFile(keys)
	if err != nil {
		return nil, err
	}

	privDer, _ := pem.Decode(privPem)
	privKey, err := x509.ParsePKCS1PrivateKey(privDer.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func (pc *PluginConfig) keysPath(path string) (string, error) {
	return filepath.Abs(filepath.Join(*pc.SignerDirectory, strings.Replace(path, "@", "_at_", -1)+".pem"))
}
