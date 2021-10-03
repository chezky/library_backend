package mail

import (
	"fmt"
	"github.com/chezky/library/db"
	"github.com/jordan-wright/email"
	"os"
	"time"

	"github.com/robfig/cron"
)

var (
	password string
	from     = "ayalabrary@gmail.com"
)

// Start gets the email password, and starts the daily cron job
func Start() {
	password = os.Getenv("EMAIL_PASS")

	// cron job for sending out emails
	c := cron.New()
	c.AddFunc("@daily", func() {
		sendBooksEmail(db.FindEmailAccounts())
	})
	c.Start()

	sendBooksEmail(db.FindEmailAccounts())
}

func sendBooksEmail(accounts []db.Account) {
	for i, acc := range accounts {
		e := email.NewEmail()

		e.From = fmt.Sprintf("AyalaBrary <%s>", from)
		e.To = []string{acc.Email}
		e.Subject = "Library Books Due Soon"
		e.HTML = []byte(createBooksEmail(acc.BookCount, arrangeBooksEmail(acc)))

		//e.Send("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"))
		fmt.Println(fmt.Sprintf("Email #%d has been sent to %s with a count of %d books", i+1, acc.Email, acc.BookCount))
	}
}


// arrangeBooksEmail creates the list of html tables that contain the books title and author, along with the due date
func arrangeBooksEmail(account db.Account) string {
	var email string

	for i, book := range account.Books {
		due := 0

		if (book.TimeStamp + 1814400) > time.Now().Unix() {
			left := book.TimeStamp + 1814400 - time.Now().Unix()
			due = time.Unix(left, 0).Day()
		}

		email += createTable(i+1, due, book.Title, book.Author)
	}
	return email
}

// createTable creates an html table that formats the title author and due date
func createTable(idx int, due int, title, author string) string {
	return `
<table border="0" cellpadding="0" cellspacing="0" align="center"
   width="100%" role="module" data-type="columns"
   style="padding:0px 0px 0px 0px;" bgcolor="#FFFFFF"
   data-distribution="1,1">
<tbody>
<tr role="module-content">
	<td height="100%" valign="top">
		<table width="290"
			   style="width:290px; border-spacing:0; border-collapse:collapse; margin:0px 10px 0px 0px;"
			   cellpadding="0" cellspacing="0" align="left"
			   border="0" bgcolor=""
			   class="column column-0">
			<tbody>
			<tr>
				<td style="padding:0px;margin:0px;border-spacing:0;">
					<table class="module" role="module"
						   data-type="text" border="0"
						   cellpadding="0" cellspacing="0"
						   width="100%"
						   style="table-layout: fixed;"
						   data-muid="079f2dee-9a7b-445c-8d17-e8f37c7a308a.1">
						<tbody>
						<tr>
							<td style="padding:18px 0px 0px 0px; line-height:22px; text-align:inherit;"
								height="100%" valign="top"
								bgcolor=""
								role="module-content">
								<div>
									<div style="font-family: inherit">
									` + fmt.Sprintf("%d", idx) + `. <span
										style="font-size: 16px"><strong>` + fmt.Sprintf("%s", title) + `</strong></span>
									</div>
									<div style="font-family: inherit">
										<span style="font-size: 16px">&nbsp;&nbsp;&nbsp;` + fmt.Sprintf("%s", author) + `</span>
									</div>
									<div></div>
								</div>
							</td>
						</tr>
						</tbody>
					</table>
				</td>
			</tr>
			</tbody>
			</table>
			<table width="290"
				   style="width:290px; border-spacing:0; border-collapse:collapse; margin:0px 0px 0px 10px;"
				   cellpadding="0" cellspacing="0" align="left"
				   border="0" bgcolor=""
				   class="column column-1">
				<tbody>
				<tr>
					<td style="padding:0px;margin:0px;border-spacing:0;">
						<table class="module" role="module"
							   data-type="text" border="0"
							   cellpadding="0" cellspacing="0"
							   width="100%"
							   style="table-layout: fixed;"
							   data-muid="acc99afe-4484-4158-9a50-7e9fed272930.1">
							<tbody>
							<tr>
								<td style="padding:18px 0px 18px 0px; line-height:22px; text-align:inherit;"
									height="100%" valign="top"
									bgcolor=""
									role="module-content">
									<div>
										<div style="font-family: inherit">
											<strong>DUE:</strong>
												` + fmt.Sprintf("%d", due) + ` Days
										</div>
										<div></div>
									</div>
								</td>
							</tr>
							</tbody>
						</table>
					</td>
				</tr>
				</tbody>
			</table>
		</td>
	</tr>
	</tbody>
</table>
`
}
