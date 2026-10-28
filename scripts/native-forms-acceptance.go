//go:build ignore

// Explicit opt-in fixture installer. Never uses live account or Drive data.
package main
import("archive/zip";"context";"encoding/json";"fmt";"io";"os";"path/filepath";"time";"github.com/oneworkspacestudio/one-workspace/internal/catalog";"github.com/oneworkspacestudio/one-workspace/internal/suite")
func main(){if e:=run();e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}}
func run()error{
 root:=filepath.Clean("C:/Projects/.cache/my-forms-20261028/fixture")
 if _,e:=os.Lstat(root);!os.IsNotExist(e){return fmt.Errorf("Fixture root must be absent")}
 password:=os.Getenv("WORKSPACE_TEST_PASSWORD");if len(password)<24{return fmt.Errorf("Set a generated test-only password")}
 z,e:=zip.OpenReader("I:/Projects/tooling/my-forms-20261028/fixture-bundle.zip");if e!=nil{return e};defer z.Close();f,e:=z.Open("catalog.json");if e!=nil{return e};b,e:=io.ReadAll(f);f.Close();if e!=nil{return e};cat,e:=catalog.Read(b);if e!=nil{return e}
 roaming:=filepath.Join(root,"roaming");if e=os.MkdirAll(roaming,0700);e!=nil{return e};os.Setenv("APPDATA",roaming)
 ctx,cancel:=context.WithTimeout(context.Background(),5*time.Minute);defer cancel()
 r:=suite.Request{Programs:filepath.Join(root,"Apps"),Data:filepath.Join(root,"Private"),Files:filepath.Join(root,"Files"),Name:"Forms test owner",Email:"forms-owner@example.test",Password:password,Confirm:password,Port:19300,Mode:"local",Offline:true,NoShortcuts:true,SelectedApps:[]string{"accounts","drive","forms"}}
 c,e:=suite.Install(ctx,r,&z.Reader,cat,func(v string){fmt.Println(v)});if e!=nil{return e}
 credentials,_:=json.Marshal(map[string]string{"email":r.Email,"password":password,"origin":c.Origin});if e=os.WriteFile(filepath.Join(c.Data,"fixture-credentials.json"),credentials,0600);e!=nil{return e}
 fmt.Println("Fixture installed with generated identities and empty Drive:",c.Origin);return nil
}
