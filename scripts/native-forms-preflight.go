//go:build ignore

package main
import("encoding/json";"fmt";"os";"reflect";"github.com/oneworkspacestudio/one-workspace/internal/catalog";"github.com/oneworkspacestudio/one-workspace/internal/distribution";"github.com/oneworkspacestudio/one-workspace/internal/suite")
func main(){b,e:=os.ReadFile(os.Args[1]);if e!=nil{panic(e)};var c suite.Config;if e=json.Unmarshal(b,&c);e!=nil{panic(e)};b,e=os.ReadFile(os.Args[2]);if e!=nil{panic(e)};incoming,e:=catalog.Read(b);if e!=nil{panic(e)};target,e:=distribution.SelectBundleUpdates(c.Catalog,incoming,[]string{"workspace"});if e!=nil{panic(e)}
for _,a:=range c.Catalog.Apps{n,ok:=target.Find(a.ID);if !ok||!reflect.DeepEqual(a,n){panic("Existing app contract changed: "+a.ID)}}
for id,p:=range c.Catalog.Packages{if id!="manager"&&target.Packages[id]!=p{panic("Existing package changed: "+id)}}
forms,ok:=target.Find("forms");if !ok||forms.Port!=c.Port+33{panic("Forms missing or incorrect managed port")};fmt.Printf("Preservation preflight passed: %d existing apps, %d existing packages; Forms port %d; target %s\n",len(c.Catalog.Apps),len(c.Catalog.Packages),forms.Port,target.Version)}
