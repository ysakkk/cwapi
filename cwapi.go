package cwapi

import (
   "fmt"
   "os"
   "path/filepath"
   "gopkg.in/yaml.v2"
   "net/http"
   "net/url"
    "strings"
  
)

// https://developer.chatwork.com/docs/endpoints
/*
 * ~/.config/cw.yml
 *    endpoint: https://api.chatwork.com/v2
 *    apitoken: xxxxxxxxxxxxxxxx
 *    roomid: xxxxxxxxxxxxxxxx
 *    alert_apitoken: xxxxxxxxxxxxxxxx
 *    alert_roomid: xxxxxxxxxxxxxxxx
 */
type Config struct {

   // API
   EndPoint      string `yaml:"endpoint"`
   RoomID        string `yaml:"roomid"`
   APIToken      string `yaml:"apitoken"`
   AlertAPIToken string `yaml:"alert_apitoken"`
   AlertRoomID   string `yaml:"alert_roomid"`
}

func loadConfig() (Config, error) {
        var c Config
        home := os.Getenv("HOME")
        configFile := filepath.Join(home, ".config", "cw.yml")
        buf, err := os.ReadFile(configFile)
        if err != nil {
                return c, fmt.Errorf("loadConfig() os.ReadFile(): %v\n", err)
        }
        if err := yaml.Unmarshal(buf, &c); err != nil {
                return c, fmt.Errorf("loadConfig() yaml.Unmarshal(): %v\n", err)
        }
        return c, nil
}


func SendAlert(body string) error {
   c, err := loadConfig() 
   if err != nil {
      return err
   }
   endpoint := c.EndPoint + "/rooms/" + c.AlertRoomID + "/messages"
   apitoken := c.AlertAPIToken
   return  sendRequest(endpoint, apitoken, body)
}


func Send(body string) error {
   c, err := loadConfig() 
   if err != nil {
      return err
   }
   endpoint := c.EndPoint + "/rooms/" + c.RoomID + "/messages"
   apitoken := c.APIToken
   return sendRequest(endpoint, apitoken, body)
}


func sendRequest(endpoint, apitoken, body string) error {

   v := url.Values{}
   v.Add("body", body)
   req, err := http.NewRequest("POST", endpoint, strings.NewReader(v.Encode()))
   if err != nil {
       return fmt.Errorf("cwapi.send() http.NewRequest(): %v\n", err)
   }

   req.Header.Set("X-ChatWorkToken", apitoken)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   client := http.Client{}

   resp, err := client.Do(req)
   if err != nil {
       return fmt.Errorf("cwapi.send() client.Do(): %v\n", err)
   }
   defer resp.Body.Close()

   return nil
}


