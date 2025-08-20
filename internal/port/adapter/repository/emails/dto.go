package port

import (
	"time"
)

// Data

type CreateFolderData struct {
	ParentId    *uint
	Name        string
	Description string
	SystemFlag  bool
}

type FilterFoldersData struct {
	Id         *[]uint
	ParentId   *[]*uint
	Name       *[]string
	SystemFlag *bool
}

type CreateEmailData struct {
	FolderId    *uint
	FromEmail   string
	FromName    string
	Subject     string
	Html        string
	Text        string
	Description string
	SystemFlag  bool
}
type FilterEmailsData struct {
	Id         *[]uint
	FolderId   *[]*uint
	SystemFlag *bool
}

type SendData struct {
	FromEmail string
	FromName  string
	Subject   string
	ToEmail   string
	Html      string
	Text      string
}

type FilterEmailLogsData struct {
	Id        *[]uint
	FromEmail *[]string
	FromName  *[]string
	ToEmail   *[]string
	Status    *[]string
	MessageId *[]string
}

// Results

type FolderResult struct {
	Id          uint
	ParentId    *uint
	Name        string
	Description string
	SystemFlag  bool
	Updated     time.Time
	Created     time.Time
}

type EmailResult struct {
	Id          uint
	FolderId    *uint
	FromEmail   string
	FromName    string
	Subject     string
	Html        string
	Text        string
	Description string
	SystemFlag  bool
	Updated     time.Time
	Created     time.Time
}

type EmailLogResult struct {
	Id        uint
	FromEmail string
	FromName  string
	Subject   string
	ToEmail   string
	Html      string
	Text      string
	Status    string
	MessageId *string
	Errors    *string
	Created   time.Time
}
