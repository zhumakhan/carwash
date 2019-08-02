package utils

import (
    c "carwashes/constants"
    "gopkg.in/gomail.v2"
    "os"
    "os/exec"
    "text/template"
    "bytes"
)

// struct to parse in mail mailTemplate
// use Vars as a handler to variables
type mailT struct {
    mail string
    subj string
    Txt string
    Vars map[string]string
}

const sendmail = "/usr/sbin/sendmail"

// just copy+paste+edit to add new mailing type
func MailWelcome(mail, name string) error {
    ma := mailT {
        mail: mail,
        subj: "Welcome!",
        Vars: map[string]string {
            "Name": name,
        },
    }
    path := c.Root + string(os.PathSeparator) + "mail_templates/welcome.txt"
    return mailTemplate(path, ma)
}

func mailTemplate(path string, ma mailT) error {
    tmpl, err := template.ParseFiles(path)
    if err != nil {
        return err
    }

    buf := &bytes.Buffer{}
    err = tmpl.Execute(buf, ma)
    if err != nil {
        return err
    }

    ma.Txt = buf.String()

    return sendRaw(ma)
}

func sendRaw(ma mailT) error {
    m := gomail.NewMessage()
    m.SetHeader("From", "turaQshare")
    m.SetHeader("To", ma.mail)
//    m.SetAddressHeader("Cc", "dan@example.com", "Dan")
    m.SetHeader("Subject", ma.subj)
    m.SetBody("text/html", ma.Txt)
//    m.Attach("/home/Alex/lolcat.jpg")

    cmd := exec.Command(sendmail, "-t")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    pw, err := cmd.StdinPipe()
    if err != nil {
        return err
    }

    err = cmd.Start()
    if err != nil {
        return err
    }

    var errs [3]error
    _, errs[0] = m.WriteTo(pw)
    errs[1] = pw.Close()
    errs[2] = cmd.Wait()
    for _, err = range errs {
        if err != nil {
            return err
        }
    }
    return nil
}
