package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"server/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// const connectionString = "mongodb+srv://admin:bH7NWTAcZWAGuKuQ@cloudassessmentsbx.dhq3mol.mongodb.net/?retryWrites=true&w=majority"

// var connectionString = "mongodb://" + "persistentcma-mongo-NLB-6ec0965dd7647173.elb.us-east-1.amazonaws.com" + ":27017"
// var connectionString = "mongodb://" + os.Getenv("MONGODB_IP") + ":27017"
//var connectionString = "mongodb://localhost:27017"
connectionString = "mongodb://cosmosdb-mizuho-eu-01:1uTuBuT8NibVNurOPw2fJRxVujOwwQtyG1LjtmbIR1JEVuAEYbC8hHuz3TqhMtp29nZvXzFXM2uqACDb7AWk0Q==@cosmosdb-mizuho-eu-01.mongo.cosmos.azure.com:10255/mizuhodevdb01?ssl=true&retryWrites=false"

// Database Name
// const dbName = "ABCSBX"
const dbName = "mizuho_db"

// Collection name
const buColl = "Business_Unit"
const agColl = "Application_Group"
const apColl = "Application"
const logColl = "logs"
const prColl = "Parameters"
const usrColl = "users"
const certColl = "Certify"

// collection object/instance
var collectionBU, collectionAG, collectionAP, collectionLG, collectionPR, collectionUsers, collectionCert *mongo.Collection
var bucket *gridfs.Bucket

type ErrorType struct {
	Message string  `json:"message"`
	Type    *string `json:"type"`
	Code    *string `json:"code"`
}
type ErrorResult struct {
	Error ErrorType `json:"error"`
}
type IDType struct {
	id string
}

type FileType struct {
	Id       string
	Filename string
}

func init() {
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collectionLG = client.Database(dbName).Collection(logColl)
	AddLog("DB initialized")

	collectionUsers = client.Database(dbName).Collection(usrColl)

	collectionBU = client.Database(dbName).Collection(buColl)
	collectionAG = client.Database(dbName).Collection(agColl)
	collectionAP = client.Database(dbName).Collection(apColl)
	collectionPR = client.Database(dbName).Collection(prColl)
	collectionCert = client.Database(dbName).Collection(certColl) // Create a new GridFS bucket
	bucket, _ = gridfs.NewBucket(collectionAP.Database())
	fmt.Println("Collections created!")
	count, err := collectionPR.EstimatedDocumentCount(context.Background())
	if err != nil || count == 0 {
		loadInitialParams()
	}
	AddLog("Collections created")
}

func loadInitialParams() {
	var Parameter = models.Parameter{
		CloudFitmentParams: []models.CloudFitmentParam{
			{
				Name:   "Technology Stack",
				Weight: 25,
				Value:  0,
				Fields: []models.CloudFitmentParam{
					{
						Name:   "Application Programming Language",
						Weight: 20,
						Value:  0,
						Fields: []models.CloudFitmentParam{
							{
								Name:  "Application_Programming_Language",
								Value: 0,
								Refs: map[string]interface{}{
									"ASP.NET":           10,
									"Bash":              0,
									"C":                 10,
									"C#":                10,
									"C++":               10,
									"Clojure":           10,
									"COBOL":             0,
									"Corba":             0,
									"F#":                10,
									"Fortran":           0,
									"Go":                10,
									"Hack":              10,
									"Haskell":           10,
									"Java":              10,
									"Javascript":        10,
									"Java, Javascript":  10,
									"Lisp":              10,
									"Perl":              10,
									"PHP":               10,
									"Powershell":        0,
									"Python":            10,
									"Ruby":              10,
									"Scala":             10,
									"Shell Script":      0,
									"Tcl":               10,
									"Visual Basic":      0,
									"Visual Basic .NET": 0,
								},
							},
						},
					},
					{
						Name:   "Web Layer Tech Stack",
						Weight: 20,
						Value:  0,
						Fields: []models.CloudFitmentParam{
							{
								Name:  "WebInterface_Layer_Technology",
								Value: 0,
								Refs: map[string]interface{}{
									"ASP.NET":                        10,
									"ASP Classic":                    4,
									"Java / J2EE / JSP / Servlets":   10,
									"HTML / CSS / Javascript":        10,
									"PHP":                            10,
									"Flash / ColdFusion":             10,
									"Citrix / VDI Technologies":      10,
									"Windows Desktop":                4,
									"None (Batch Process, IDE, etc)": 10,
									"None (Unix App, Shell / Terminal Application)": 0,
								},
							},
						},
					},
					{
						Name:   "Web Layer OS",
						Weight: 20,
						Value:  0,
						Fields: []models.CloudFitmentParam{
							{
								Name:  "WebInterface_Layer_Operating_System",
								Value: 0,
								Refs: map[string]interface{}{
									"None":                                           10,
									"Custom Linux Distro":                            0,
									"Windows":                                        10,
									"Linux":                                          10,
									"Windows 2000":                                   4,
									"Windows 2003":                                   4,
									"Windows 2008":                                   4,
									"Windows 2008 R2":                                7,
									"Windows 2012":                                   10,
									"Windows 2016":                                   10,
									"Windows 7":                                      7,
									"Windows 8":                                      10,
									"Windows 10":                                     10,
									"Red Hat Enterprise Linux (RHEL) 2.x":            4,
									"Red Hat Enterprise Linux (RHEL) 3.x":            4,
									"Red Hat Enterprise Linux (RHEL) 4.x":            4,
									"Red Hat Enterprise Linux (RHEL) 5.x":            4,
									"Red Hat Enterprise Linux (RHEL) 6.x":            10,
									"Red Hat Enterprise Linux (RHEL) 7.x":            10,
									"SUSE Linux Enterprise Server (SLES) 8.x":        4,
									"SUSE Linux Enterprise Server (SLES) 9.x":        4,
									"SUSE Linux Enterprise Server (SLES) 10.x":       4,
									"SUSE Linux Enterprise Server (SLES) 11.x":       7,
									"SUSE Linux Enterprise Server (SLES) 12.x":       10,
									"CentOS 3.x":                                     4,
									"CentOS 4.x":                                     4,
									"CentOS 5.x":                                     4,
									"CentOS 6.x":                                     7,
									"CentOS 7.x":                                     10,
									"Ubuntu Linux 4.x":                               4,
									"Ubuntu Linux 5.x":                               4,
									"Ubuntu Linux 6.x":                               4,
									"Ubuntu Linux 7.x":                               4,
									"Ubuntu Linux 8.x":                               4,
									"Ubuntu Linux 9.x":                               4,
									"Ubuntu Linux 10.x":                              4,
									"Ubuntu Linux 11.x":                              4,
									"Ubuntu Linux 12.x":                              4,
									"Ubuntu Linux 13.x":                              4,
									"Ubuntu Linux 14.x":                              10,
									"Ubuntu Linux 15.x":                              10,
									"Ubuntu Linux 16.x":                              10,
									"Ubuntu Linux 17.x":                              10,
									"Linux_Oracle":                                   7,
									"Linux_Debian":                                   7,
									"VMware ESX/ESXi 4.x":                            4,
									"VMware ESX/ESXi 5.x":                            4,
									"VMware ESX/ESXi 6.x":                            4,
									"CoreOS":                                         4,
									"Solaris 8":                                      0,
									"Solaris 9":                                      0,
									"Solaris 10 (SPARC)":                             0,
									"None_strat":                                     "Rehost",
									"Custom Linux Distro_strat":                      "Replatform/Refactor/Rearchitect",
									"Windows_strat":                                  "Rehost",
									"Linux_strat":                                    "Rehost",
									"Windows 2000_strat":                             "Replatform",
									"Windows 2003_strat":                             "Replatform",
									"Windows 2008_strat":                             "Replatform",
									"Windows 2008 R2_strat":                          "Rehost/Replatform",
									"Windows 2012_strat":                             "Rehost",
									"Windows 2016_strat":                             "Rehost",
									"Windows 7_strat":                                "Rehost/Replatform",
									"Windows 8_strat":                                "Rehost",
									"Windows 10_strat":                               "Rehost",
									"Red Hat Enterprise Linux (RHEL) 2.x_strat":      "Replatform",
									"Red Hat Enterprise Linux (RHEL) 3.x_strat":      "Replatform",
									"Red Hat Enterprise Linux (RHEL) 4.x_strat":      "Replatform",
									"Red Hat Enterprise Linux (RHEL) 5.x_strat":      "Replatform",
									"Red Hat Enterprise Linux (RHEL) 6.x_strat":      "Rehost",
									"Red Hat Enterprise Linux (RHEL) 7.x_strat":      "Rehost",
									"SUSE Linux Enterprise Server (SLES) 8.x_strat":  "Replatform",
									"SUSE Linux Enterprise Server (SLES) 9.x_strat":  "Replatform",
									"SUSE Linux Enterprise Server (SLES) 10.x_strat": "Replatform",
									"SUSE Linux Enterprise Server (SLES) 11.x_strat": "Rehost/Replatform",
									"SUSE Linux Enterprise Server (SLES) 12.x_strat": "Rehost",
									"CentOS 3.x_strat":                               "Replatform",
									"CentOS 4.x_strat":                               "Replatform",
									"CentOS 5.x_strat":                               "Replatform",
									"CentOS 6.x_strat":                               "Rehost/Replatform",
									"CentOS 7.x_strat":                               "Rehost",
									"Ubuntu Linux 4.x_strat":                         "Replatform",
									"Ubuntu Linux 5.x_strat":                         "Replatform",
									"Ubuntu Linux 6.x_strat":                         "Replatform",
									"Ubuntu Linux 7.x_strat":                         "Replatform",
									"Ubuntu Linux 8.x_strat":                         "Replatform",
									"Ubuntu Linux 9.x_strat":                         "Replatform",
									"Ubuntu Linux 10.x_strat":                        "Replatform",
									"Ubuntu Linux 11.x_strat":                        "Replatform",
									"Ubuntu Linux 12.x_strat":                        "Replatform",
									"Ubuntu Linux 13.x_strat":                        "Replatform",
									"Ubuntu Linux 14.x_strat":                        "Rehost",
									"Ubuntu Linux 15.x_strat":                        "Rehost",
									"Ubuntu Linux 16.x_strat":                        "Rehost",
									"Ubuntu Linux 17.x_strat":                        "Rehost",
									"Linux_Oracle_strat":                             "Rehost/Replatform",
									"Linux_Debian_strat":                             "Rehost/Replatform",
									"VMware ESX/ESXi 4.x_strat":                      "Replatform",
									"VMware ESX/ESXi 5.x_strat":                      "Replatform",
									"VMware ESX/ESXi 6.x_strat":                      "Replatform",
									"CoreOS_strat":                                   "Replatform",
									"Solaris 8_strat":                                "Rearchitect",
									"Solaris 9_strat":                                "Rearchitect",
									"Solaris 10 (SPARC)_strat":                       "Rearchitect",
									"Solaris 10 (x86/x64)_strat":                     "Refactor/Rearchitect or Host on Citrix Cloud",
									"Solaris 11_strat":                               "Refactor/Rearchitect or Host on Oracle Cloud",
									"AIX 5.1_strat":                                  "Rearchitect",
									"AIX 5.2_strat":                                  "Refactor/Rearchitect or Host on Citrix Cloud",
									"AIX 5.3_strat":                                  "Refactor/Rearchitect or Host on Citrix Cloud",
									"HP-UX 11_strat":                                 "Refactor/Rearchitect or Host on Citrix Cloud",
									"HP-UX 11i_strat":                                "Refactor/Rearchitect or Host on Citrix Cloud",
									"Unix_strat":                                     "Refactor/Rearchitect or Host on Citrix Cloud",
									"AS400_strat":                                    "Rearchitect",
									"OS/2_strat":                                     "Rearchitect",
									"z/OS_strat":                                     "Rearchitect",
								},
							},
						},
					},
					{
						Name:       "Middleware Tech Stack",
						Depends:    "Has_Middleware",
						DependsVal: "Yes",
						EmptyVal:   10,
						Weight:     20,
						Value:      0,
						Fields: []models.CloudFitmentParam{
							{
								Name:  "Middleware_Technology",
								Value: 0,
								Refs: map[string]interface{}{
									"None":                                      10,
									"BEA Tuxedo":                                0,
									"BEA Weblogic Family":                       0,
									"Borland Application Server":                0,
									"FioranoMQ":                                 10,
									"HP Process Manager":                        0,
									"HP/Bluestone Total-e-Server":               0,
									"HP/Bluestone Total-e-Transactions":         0,
									"IBM CICS, IMS, TPF, TxSeries":              0,
									"IBM MQ":                                    10,
									"IBM MQ Integrator":                         0,
									"IBM MQ Workflow":                           0,
									"IBM WebMethods":                            0,
									"IBM WebSphere Family":                      0,
									"iPlanet Application Server":                0,
									"iPlanet Process Builder":                   0,
									"Apache Tomcat":                             0,
									"JBoss Server":                              0,
									"Mercator Enterprise Broker":                0,
									"Microsoft BizTalk Server":                  0,
									"Microsoft MQ":                              10,
									"MS-RPC":                                    10,
									"JDBC libraries":                            10,
									"ODBC libraries":                            10,
									"SAP XI / PI Family":                        0,
									"SeeBeyond eExchange Integration Suite":     0,
									"Software AG EntireX":                       0,
									"SonicMQ":                                   10,
									"Sun Java MQ":                               10,
									"Sun NFS":                                   10,
									"SwiftMQ":                                   10,
									"Tibco ActiveEnterprise ActivePortal":       0,
									"Tibco Rendezvous":                          10,
									"Tibco Software ActiveEnterprise InConcert": 0,
									"Tibco Software ActiveEnterprise Integration Manager": 0,
									"Vitria BusinessWare Automator":                       0,
									"Windows NFS":                                         10,
									"XML-Server":                                          0,
								},
							},
						},
					},
					{
						Name:       "Middleware OS",
						Depends:    "Has_Middleware",
						DependsVal: "Yes",
						EmptyVal:   10,
						Weight:     20,
						Value:      0,
						Fields: []models.CloudFitmentParam{
							{
								Name:  "Middleware_Operating_System",
								Value: 0,
								Refs: map[string]interface{}{
									"None":                                           10,
									"Custom Linux Distro":                            0,
									"Windows":                                        10,
									"Linux":                                          10,
									"Windows 2000":                                   4,
									"Windows 2003":                                   4,
									"Windows 2008":                                   4,
									"Windows 2008 R2":                                7,
									"Windows 2012":                                   10,
									"Windows 2016":                                   10,
									"Windows 7":                                      7,
									"Windows 8":                                      10,
									"Windows 10":                                     10,
									"Red Hat Enterprise Linux (RHEL) 2.x":            4,
									"Red Hat Enterprise Linux (RHEL) 3.x":            4,
									"Red Hat Enterprise Linux (RHEL) 4.x":            4,
									"Red Hat Enterprise Linux (RHEL) 5.x":            4,
									"Red Hat Enterprise Linux (RHEL) 6.x":            10,
									"Red Hat Enterprise Linux (RHEL) 7.x":            10,
									"SUSE Linux Enterprise Server (SLES) 8.x":        4,
									"SUSE Linux Enterprise Server (SLES) 9.x":        4,
									"SUSE Linux Enterprise Server (SLES) 10.x":       4,
									"SUSE Linux Enterprise Server (SLES) 11.x":       7,
									"SUSE Linux Enterprise Server (SLES) 12.x":       10,
									"CentOS 3.x":                                     4,
									"CentOS 4.x":                                     4,
									"CentOS 5.x":                                     4,
									"CentOS 6.x":                                     7,
									"CentOS 7.x":                                     10,
									"Ubuntu Linux 4.x":                               4,
									"Ubuntu Linux 5.x":                               4,
									"Ubuntu Linux 6.x":                               4,
									"Ubuntu Linux 7.x":                               4,
									"Ubuntu Linux 8.x":                               4,
									"Ubuntu Linux 9.x":                               4,
									"Ubuntu Linux 10.x":                              4,
									"Ubuntu Linux 11.x":                              4,
									"Ubuntu Linux 12.x":                              4,
									"Ubuntu Linux 13.x":                              4,
									"Ubuntu Linux 14.x":                              10,
									"Ubuntu Linux 15.x":                              10,
									"Ubuntu Linux 16.x":                              10,
									"Ubuntu Linux 17.x":                              10,
									"Linux_Oracle":                                   7,
									"Linux_Debian":                                   7,
									"VMware ESX/ESXi 4.x":                            4,
									"VMware ESX/ESXi 5.x":                            4,
									"VMware ESX/ESXi 6.x":                            4,
									"CoreOS":                                         4,
									"Solaris 8":                                      0,
									"Solaris 9":                                      0,
									"Solaris 10 (SPARC)":                             0,
									"None_strat":                                     "Rehost",
									"Custom Linux Distro_strat":                      "Replatform/Refactor/Rearchitect",
									"Windows_strat":                                  "Rehost",
									"Linux_strat":                                    "Rehost",
									"Windows 2000_strat":                             "Replatform",
									"Windows 2003_strat":                             "Replatform",
									"Windows 2008_strat":                             "Replatform",
									"Windows 2008 R2_strat":                          "Rehost/Replatform",
									"Windows 2012_strat":                             "Rehost",
									"Windows 2016_strat":                             "Rehost",
									"Windows 7_strat":                                "Rehost/Replatform",
									"Windows 8_strat":                                "Rehost",
									"Windows 10_strat":                               "Rehost",
									"Red Hat Enterprise Linux (RHEL) 2.x_strat":      "Replatform",
									"Red Hat Enterprise Linux (RHEL) 3.x_strat":      "Replatform",
									"Red Hat Enterprise Linux (RHEL) 4.x_strat":      "Replatform",
									"Red Hat Enterprise Linux (RHEL) 5.x_strat":      "Replatform",
									"Red Hat Enterprise Linux (RHEL) 6.x_strat":      "Rehost",
									"Red Hat Enterprise Linux (RHEL) 7.x_strat":      "Rehost",
									"SUSE Linux Enterprise Server (SLES) 8.x_strat":  "Replatform",
									"SUSE Linux Enterprise Server (SLES) 9.x_strat":  "Replatform",
									"SUSE Linux Enterprise Server (SLES) 10.x_strat": "Replatform",
									"SUSE Linux Enterprise Server (SLES) 11.x_strat": "Rehost/Replatform",
									"SUSE Linux Enterprise Server (SLES) 12.x_strat": "Rehost",
									"CentOS 3.x_strat":                               "Replatform",
									"CentOS 4.x_strat":                               "Replatform",
									"CentOS 5.x_strat":                               "Replatform",
									"CentOS 6.x_strat":                               "Rehost/Replatform",
									"CentOS 7.x_strat":                               "Rehost",
									"Ubuntu Linux 4.x_strat":                         "Replatform",
									"Ubuntu Linux 5.x_strat":                         "Replatform",
									"Ubuntu Linux 6.x_strat":                         "Replatform",
									"Ubuntu Linux 7.x_strat":                         "Replatform",
									"Ubuntu Linux 8.x_strat":                         "Replatform",
									"Ubuntu Linux 9.x_strat":                         "Replatform",
									"Ubuntu Linux 10.x_strat":                        "Replatform",
									"Ubuntu Linux 11.x_strat":                        "Replatform",
									"Ubuntu Linux 12.x_strat":                        "Replatform",
									"Ubuntu Linux 13.x_strat":                        "Replatform",
									"Ubuntu Linux 14.x_strat":                        "Rehost",
									"Ubuntu Linux 15.x_strat":                        "Rehost",
									"Ubuntu Linux 16.x_strat":                        "Rehost",
									"Ubuntu Linux 17.x_strat":                        "Rehost",
									"Linux_Oracle_strat":                             "Rehost/Replatform",
									"Linux_Debian_strat":                             "Rehost/Replatform",
									"VMware ESX/ESXi 4.x_strat":                      "Replatform",
									"VMware ESX/ESXi 5.x_strat":                      "Replatform",
									"VMware ESX/ESXi 6.x_strat":                      "Replatform",
									"CoreOS_strat":                                   "Replatform",
									"Solaris 8_strat":                                "Rearchitect",
									"Solaris 9_strat":                                "Rearchitect",
									"Solaris 10 (SPARC)_strat":                       "Rearchitect",
									"Solaris 10 (x86/x64)_strat":                     "Refactor/Rearchitect or Host on Citrix Cloud",
									"Solaris 11_strat":                               "Refactor/Rearchitect or Host on Oracle Cloud",
									"AIX 5.1_strat":                                  "Rearchitect",
									"AIX 5.2_strat":                                  "Refactor/Rearchitect or Host on Citrix Cloud",
									"AIX 5.3_strat":                                  "Refactor/Rearchitect or Host on Citrix Cloud",
									"HP-UX 11_strat":                                 "Refactor/Rearchitect or Host on Citrix Cloud",
									"HP-UX 11i_strat":                                "Refactor/Rearchitect or Host on Citrix Cloud",
									"Unix_strat":                                     "Refactor/Rearchitect or Host on Citrix Cloud",
									"AS400_strat":                                    "Rearchitect",
									"OS/2_strat":                                     "Rearchitect",
									"z/OS_strat":                                     "Rearchitect",
								},
							},
						},
					},
				},
			},
			{
				Name:   "Architecture Dependencies",
				Weight: 25,
				Value:  0,
				Fields: []models.CloudFitmentParam{
					{
						Name:   "CPU Architecture",
						Weight: 25,
						Value:  0,
						Fields: []models.CloudFitmentParam{
							{
								Name:  "CPU_Architecture",
								Value: 0,
								Refs: map[string]interface{}{
									"x86-64 bit":   10,
									"x86-32 bit":   5,
									"Oracle SPARC": 0,
									"IBM Z":        0,
									"IBM Power":    0,
									"HP Itanium":   0,
									"ARM":          0,
								},
							},
						},
					},
					{
						Name:   "Hardware / OS Dependence",
						Weight: 25,
						Value:  0,
						Fields: []models.CloudFitmentParam{
							{
								Name:  "Hardware_OS_Dependency",
								Value: 0,
								Refs: map[string]interface{}{
									"No":  10,
									"Yes": 0,
								},
							},
						},
					},
					{
						Name:   "Use of Shared Disks",
						Weight: 25,
						Value:  0,
						Fields: []models.CloudFitmentParam{
							{
								Name:  "Shared_Disk_Architecture",
								Value: 0,
								Refs: map[string]interface{}{
									"No":  10,
									"Yes": 0,
								},
							},
						},
					},
					{
						Name:   "Use of IP Multicast",
						Weight: 25,
						Value:  0,
						Fields: []models.CloudFitmentParam{
							{
								Name:  "IP_Multicast",
								Value: 0,
								Refs: map[string]interface{}{
									"No":  10,
									"Yes": 0,
								},
							},
						},
					},
				},
			},
			{
				Name:   "Regulatory & Compliance",
				Weight: 25,
				Value:  0,
				Fields: []models.CloudFitmentParam{
					{
						Name:       "Compliance Requirement Fitment",
						Weight:     100,
						Depends:    "Has_Compliance_Requirement",
						DependsVal: "Yes",
						EmptyVal:   10,
						Value:      0,
						Fields: []models.CloudFitmentParam{
							{
								Name:  "Compliance_Standards",
								Value: 0,
								Refs: map[string]interface{}{
									"Not listed below":             0,
									"C5 (Germany)":                 10,
									"Cyber Essentials Plus (UK)":   10,
									"DoD SRG":                      10,
									"FedRAMP":                      10,
									"FIPS":                         10,
									"IRAP (Australia)":             10,
									"ISO 9001":                     10,
									"ISO 27001":                    10,
									"ISO 27017":                    10,
									"ISO 27018":                    10,
									"MTCS (Singapore)":             10,
									"PCI DSS ":                     10,
									"SEC Rule 17-a-4(F)":           10,
									"SOC 1":                        10,
									"SOC 2":                        10,
									"SOC 3":                        10,
									"CISPE":                        10,
									"FERPA":                        10,
									"GLBA":                         10,
									"HIPAA":                        10,
									"HITECH":                       10,
									"IRS 1075":                     10,
									"ITAR":                         10,
									"My Number Act (Japan)":        10,
									"U.K. DPA 1988":                10,
									"VPAT / Section 508":           10,
									"EU Data Protection Directive": 10,
									"Privacy Act (Australia)":      10,
									"Privacy Act (New Zeland)":     10,
									"PDPA 2010 (Malaysia)":         10,
									"PDPA 2012 (Singapore)":        10,
									"PIPEDA (Canada)":              10,
									"Spanish DPA Authorization":    10,
									"CIS":                          10,
									"CJIS":                         10,
									"CSA":                          10,
									"ENS (Spain)":                  10,
								},
							},
						},
					},
				},
			},
			{
				Name:   "Licensing",
				Weight: 25,
				Value:  0,
				Fields: []models.CloudFitmentParam{
					{
						Name:       "License Portability",
						Weight:     100,
						Depends:    "Has_Licensing_Requirement",
						DependsVal: "Yes",
						EmptyVal:   10,
						Value:      0,
						Fields: []models.CloudFitmentParam{
							{
								Name:  "License_Portability",
								Value: 0,
								Refs: map[string]interface{}{
									"License not portable, Cloud license not available": 0,
									"License not portable, Cloud license available":     10,
									"License portable": 10,
								},
							},
						},
					},
				},
			},
		},
	}
	insertResult, err := collectionPR.InsertOne(context.Background(), Parameter)
	if err != nil {
		AddLog("Unable to load initial parameters, " + err.Error())
	} else if insertResult != nil {
		AddLog("Loaded initial parameters")
	} else {
		AddLog("Unable to load initial parameters")
	}
}

func AddLog(log string) {
	var LogTypeColl models.LogType
	LogTypeColl.Created = time.Now()
	LogTypeColl.Log = log
	insertResult, err := collectionLG.InsertOne(context.Background(), LogTypeColl)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(insertResult)
	}
}
func GetAllParameters(w http.ResponseWriter, r *http.Request) {
	payload, err := getAllParameters()
	if err != nil {
		AddLog("GetAllParameters failed: err:" + err.Error())
		var resp ErrorResult
		resp.Error.Message = err.Error()
		json.NewEncoder(w).Encode(resp)
	} else if payload == nil {
		AddLog("GetAllParameters failed: err:404")
		w.WriteHeader(http.StatusNotFound)
	} else {
		AddLog("GetAllParameters success")
		json.NewEncoder(w).Encode(payload)
	}
}
func getAllParameters() (primitive.M, error) {
	var result primitive.M
	err := collectionPR.FindOne(context.Background(), bson.M{}).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, err
}

// #region Business Unit
func GetAllBusinessUnit(w http.ResponseWriter, r *http.Request) {
	towerLeads, err := getAllTowerLeads()
	if err != nil {
		AddLog("GetAllTowerLeads failed: err: " + err.Error())
		var resp ErrorResult
		resp.Error.Message = err.Error()
		json.NewEncoder(w).Encode(resp)
	} else {
		AddLog("GetAllTowerLeads success")
		json.NewEncoder(w).Encode(towerLeads)
	}
}
func GetBusinessUnit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	if id == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	payload, err := getOneBusinessUnit(id)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else if payload == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(payload)
	}
}
func CreateBusinessUnit(w http.ResponseWriter, r *http.Request) {
	var BusinessUnitColl models.BusinessUnitType
	_ = json.NewDecoder(r.Body).Decode(&BusinessUnitColl)
	if BusinessUnitColl.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing Name"
		json.NewEncoder(w).Encode(resp)
		return
	}
	BusinessUnitColl.Updated = time.Now()
	BusinessUnitColl.Type = "Business_Unit"
	insertResult, err := insertOneBusinessUnit(BusinessUnitColl)
	if err != nil {
		AddLog("CreateBusinessUnit failed: err:" + err.Error())
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		AddLog("GetAllBusinessUnit success")
		json.NewEncoder(w).Encode(insertResult)
	}
}
func UpdateBusinessUnit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	if id == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	var BusinessUnitColl models.BusinessUnitType
	_ = json.NewDecoder(r.Body).Decode(&BusinessUnitColl)
	if BusinessUnitColl.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing Name"
		json.NewEncoder(w).Encode(resp)
		return
	}
	BusinessUnitColl.ID = id
	BusinessUnitColl.Updated = time.Now()
	BusinessUnitColl.Type = "Business_Unit"
	updateResult, err := updateOneBusinessUnit(BusinessUnitColl.ID, BusinessUnitColl)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else if updateResult.MatchedCount == 0 {
		var resp ErrorResult
		resp.Error.Message = "Record not found"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
	} else {
		json.NewEncoder(w).Encode(updateResult)
	}
}
func DeleteBusinessUnit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	if id == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	deleteResult, err := deleteOneBusinessUnit(id)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		json.NewEncoder(w).Encode(deleteResult)
	}
}
func getAllTowerLeads() ([]string, error) {
	var towerLeads []string

	// Query to get distinct Tower_Lead values from the "Application" collection
	cur, err := collectionAP.Distinct(context.Background(), "Tower_Lead", bson.M{})
	if err != nil {
		return towerLeads, err
	}

	// Assert each value to string and append to the result slice
	for _, value := range cur {
		if towerLead, ok := value.(string); ok {
			towerLeads = append(towerLeads, towerLead)
		}
	}

	return towerLeads, nil
}
func getOneBusinessUnit(id primitive.ObjectID) (primitive.M, error) {
	var result primitive.M
	err := collectionBU.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result)
	AddLog("getAllBusinessUnitRecords success")
	return result, err
}
func insertOneBusinessUnit(BusinessUnit models.BusinessUnitType) (*mongo.InsertOneResult, error) {
	insertResult, err := collectionBU.InsertOne(context.Background(), BusinessUnit)
	return insertResult, err
}
func updateOneBusinessUnit(id primitive.ObjectID, BusinessUnit models.BusinessUnitType) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	updateResult, err := collectionBU.ReplaceOne(context.Background(), filter, BusinessUnit)
	return updateResult, err
}
func deleteOneBusinessUnit(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{
		"$and": bson.A{bson.M{"_id": id},
			bson.M{"type": "Business_Unit"},
		},
	}
	d, err := collectionBU.DeleteOne(context.Background(), filter)
	fmt.Println("delete a Single BU Record ")
	return d, err
}

// #endregion

// #region Application Group
func GetAllApplicationGroup(w http.ResponseWriter, r *http.Request) {
	supportGroups, err := getAllSupportGroups()
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		json.NewEncoder(w).Encode(resp)
	} else {
		fmt.Println("List of Support Groups:", strings.Join(supportGroups, ","))
		json.NewEncoder(w).Encode(supportGroups)
	}
}
func GetApplicationGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	if id == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	payload, err := getOneApplicationGroup(id)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else if payload == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(payload)
	}
}
func CreateApplicationGroup(w http.ResponseWriter, r *http.Request) {
	var ApplicationGroupColl models.ApplicationGroupType
	_ = json.NewDecoder(r.Body).Decode(&ApplicationGroupColl)
	if ApplicationGroupColl.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing Name"
		json.NewEncoder(w).Encode(resp)
		return
	}
	if ApplicationGroupColl.Business_Unit_ID == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing BU ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	ApplicationGroupColl.Updated = time.Now()
	ApplicationGroupColl.Type = "Application_Group"
	insertResult, err := insertOneApplicationGroup(ApplicationGroupColl)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		json.NewEncoder(w).Encode(insertResult)
	}
}
func UpdateApplicationGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var ApplicationGroupColl models.ApplicationGroupType
	_ = json.NewDecoder(r.Body).Decode(&ApplicationGroupColl)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	if ApplicationGroupColl.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing Name"
		json.NewEncoder(w).Encode(resp)
		return
	}
	if ApplicationGroupColl.Business_Unit_ID == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing BU ID"
		json.NewEncoder(w).Encode(resp)
		return
	}

	if id == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing AG ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	ApplicationGroupColl.ID = id
	ApplicationGroupColl.Updated = time.Now()
	ApplicationGroupColl.Type = "Application_Group"
	updateResult, err := updateOneApplicationGroup(ApplicationGroupColl.ID, ApplicationGroupColl)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		json.NewEncoder(w).Encode(updateResult)
	}
}
func DeleteApplicationGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	if id == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	deleteResult, err := deleteOneApplicationGroup(id)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		json.NewEncoder(w).Encode(deleteResult)
	}
}
func getAllSupportGroups() ([]string, error) {
	var supportGroups []string

	cur, err := collectionAP.Distinct(context.Background(), "Support_Group", bson.M{})
	if err != nil {
		return supportGroups, err
	}

	for _, value := range cur {
		if supportGroup, ok := value.(string); ok {
			supportGroups = append(supportGroups, supportGroup)
		}
	}

	return supportGroups, nil
}
func getOneApplicationGroup(id primitive.ObjectID) (primitive.M, error) {
	var result primitive.M
	qry := []bson.M{
		{
			"$match": bson.M{
				"type": "Application_Group",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "Business_Unit",    // Child collection to join
				"localField":   "business_unit_id", // Parent collection reference holding field
				"foreignField": "_id",              // Child collection reference field
				"as":           "Business_Unit",    // Arbitrary field name to store result set
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$Business_Unit",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}
	cur, err := collectionAG.Aggregate(context.Background(), qry)
	if err != nil {
		return result, err
	}
	if cur.Next(context.Background()) {
		e := cur.Decode(&result)
		if e != nil {
			return result, err
		}
	}
	return result, err
}
func insertOneApplicationGroup(ApplicationGroup models.ApplicationGroupType) (*mongo.InsertOneResult, error) {
	insertResult, err := collectionAG.InsertOne(context.Background(), ApplicationGroup)
	fmt.Println("insert a Single AG Record ", insertResult.InsertedID)
	return insertResult, err
}
func updateOneApplicationGroup(id primitive.ObjectID, ApplicationGroup models.ApplicationGroupType) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	updateResult, err := collectionAG.ReplaceOne(context.Background(), filter, ApplicationGroup)
	fmt.Println("update a Single AG Record ", updateResult.MatchedCount)
	return updateResult, err
}
func deleteOneApplicationGroup(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{
		"$and": bson.A{bson.M{"_id": id},
			bson.M{"type": "Application_Group"},
		},
	}
	d, err := collectionAG.DeleteOne(context.Background(), filter)
	return d, err
}

// #endregion

// #region Application
func GetAllApplication(w http.ResponseWriter, r *http.Request) {
	payload, err := getAllApplicationRecords()
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		json.NewEncoder(w).Encode(resp)
	} else if payload == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(payload)
	}
}
func GetApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	if id == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	payload, err := getOneApplication(id)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else if payload == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(payload)
	}
}
func CreateApplication(w http.ResponseWriter, r *http.Request) {
	var ApplicationColl models.ApplicationType
	_ = json.NewDecoder(r.Body).Decode(&ApplicationColl)
	if ApplicationColl.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing Name"
		json.NewEncoder(w).Encode(resp)
		return
	}
	// if ApplicationColl.Application_Group_ID == primitive.NilObjectID {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	var resp ErrorResult
	// 	resp.Error.Message = "Missing AG ID"
	// 	json.NewEncoder(w).Encode(resp)
	// 	return
	// }
	ApplicationColl.Updated = time.Now()
	ApplicationColl.Type = "Application"
	insertResult, err := insertOneApplication(ApplicationColl)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		json.NewEncoder(w).Encode(insertResult)
	}
}
func UpdateApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var ApplicationColl models.ApplicationType
	var CertificationColl models.Certification
	_ = json.NewDecoder(r.Body).Decode(&ApplicationColl)
	if ApplicationColl.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing Name"
		json.NewEncoder(w).Encode(resp)
		return
	}
	// if ApplicationColl.Application_Group_ID == primitive.NilObjectID {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	var resp ErrorResult
	// 	resp.Error.Message = "Missing AG ID"
	// 	json.NewEncoder(w).Encode(resp)
	// 	return
	// }

	if id == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing AP ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	ApplicationColl.ID = id
	ApplicationColl.Updated = time.Now()
	ApplicationColl.Type = "Application"

	if ApplicationColl.Certified {
		CertificationColl.Updated = time.Now()
		CertificationColl.Certified_by = ApplicationColl.Modified_By
		CertificationColl.CMDB_ID = ApplicationColl.CMDB_ID
		insertResult, err1 := UpdateCertification(CertificationColl)
		if err1 != nil {
			var resp ErrorResult
			resp.Error.Message = err1.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
		}
		ApplicationColl.LogIds = append(ApplicationColl.LogIds, insertResult.InsertedID.(primitive.ObjectID))
	}

	if ApplicationColl.Certified {
		ApplicationColl.Certified = false
	}

	updateResult, err := updateOneApplication(ApplicationColl.ID, ApplicationColl)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		json.NewEncoder(w).Encode(updateResult)
	}
}
func DeleteApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	if id == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	deleteResult, err := deleteOneApplication(id)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		json.NewEncoder(w).Encode(deleteResult)
	}
}
func UpdateManyApplication(w http.ResponseWriter, r *http.Request) {
	var ApplicationColl []models.ApplicationType
	_ = json.NewDecoder(r.Body).Decode(&ApplicationColl)
	// fmt.Printf("%+v\n", ApplicationColl)
	updateResult, err := updateManyApplication(ApplicationColl)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		json.NewEncoder(w).Encode(updateResult)
	}
}
func getAllApplicationRecords() ([]primitive.M, error) {
	var results []primitive.M
	qry := []bson.M{
		// {
		// 	"$match": bson.M{
		// 		"type": "Application",
		// 	},
		// },
		{
			"$lookup": bson.M{
				"from":         "Certify",
				"localField":   "logids",
				"foreignField": "_id",
				"as":           "logs",
			},
		},
		{
			"$project": bson.M{
				// "type":                 0,
				"application_group_id": 0,
			},
		},
		// },
		// {
		// 	"$unwind": bson.M{
		// 		"path":                       "$Application_Group",
		// 		"preserveNullAndEmptyArrays": true,
		// 	},
		// },
		// {
		// 	"$lookup": bson.M{
		// 		"from":         "Business_Unit",                      // Child collection to join
		// 		"localField":   "Application_Group.business_unit_id", // Parent collection reference holding field
		// 		"foreignField": "_id",                                // Child collection reference field
		// 		"as":           "Business_Unit",                      // Arbitrary field name to store result set
		// 	},
		// },
		// {
		// 	"$unwind": bson.M{
		// 		"path":                       "$Business_Unit",
		// 		"preserveNullAndEmptyArrays": true,
		// 	},
		// },
	}
	cur, err := collectionAP.Aggregate(context.Background(), qry)
	if err != nil {
		return results, err
	}
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			continue
		}
		results = append(results, result)
	}
	cur.Close(context.Background())
	return results, err
}
func getOneApplication(id primitive.ObjectID) (primitive.M, error) {
	var result primitive.M
	qry := []bson.M{
		{
			"$match": bson.M{
				// "type": "Application",
				"_id": id,
			},
		},
		// {
		// 	"$lookup": bson.M{
		// 		"from":         "Application_Group",    // Child collection to join
		// 		"localField":   "application_group_id", // Parent collection reference holding field
		// 		"foreignField": "_id",                  // Child collection reference field
		// 		"as":           "Application_Group",    // Arbitrary field name to store result set
		// 	},
		// },
		{
			"$unwind": bson.M{
				"path":                       "$Application_Group",
				"preserveNullAndEmptyArrays": true,
			},
		},
		// {
		// 	"$lookup": bson.M{
		// 		"from":         "Business_Unit",                      // Child collection to join
		// 		"localField":   "Application_Group.business_unit_id", // Parent collection reference holding field
		// 		"foreignField": "_id",                                // Child collection reference field
		// 		"as":           "Business_Unit",                      // Arbitrary field name to store result set
		// 	},
		// },
		{
			"$lookup": bson.M{
				"from":         "Certify",
				"localField":   "logids",
				"foreignField": "_id",
				"as":           "logs",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$Business_Unit",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}
	cur, err := collectionAP.Aggregate(context.Background(), qry)
	if err != nil {
		return result, err
	}
	if cur.Next(context.Background()) {
		e := cur.Decode(&result)
		if e != nil {
			return result, err
		}
	}
	return result, err
}
func insertOneApplication(Application models.ApplicationType) (*mongo.InsertOneResult, error) {
	insertResult, err := collectionAP.InsertOne(context.Background(), Application)
	fmt.Println("insert a Single AP Record ", insertResult.InsertedID)
	return insertResult, err
}
func updateOneApplication(id primitive.ObjectID, Application models.ApplicationType) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	updateResult, err := collectionAP.ReplaceOne(context.Background(), filter, Application)
	fmt.Println("update a Single AP Record ", updateResult.MatchedCount)
	return updateResult, err
}
func deleteOneApplication(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{
		"$and": bson.A{bson.M{"_id": id},
			bson.M{"type": "Application"},
		},
	}
	d, err := collectionAP.DeleteOne(context.Background(), filter)
	return d, err
}
func updateManyApplication(Applications []models.ApplicationType) (*mongo.BulkWriteResult, error) {

	models := []mongo.WriteModel{}
	for _, v := range Applications {
		updateRec := bson.D{}
		if v.Name != "" {
			updateRec = append(updateRec, bson.E{Key: "name", Value: v.Name})
		}
		// if v.Application_Group_ID != primitive.NilObjectID {
		// 	updateRec = append(updateRec, bson.E{Key: "application_group_id", Value: v.Application_Group_ID})
		// }
		// updateRec = append(updateRec, bson.E{Key: "type", Value: "Application"})
		if v.Modified_By != "" {
			updateRec = append(updateRec, bson.E{Key: "modified_by", Value: v.Modified_By})
			updateRec = append(updateRec, bson.E{Key: "updated", Value: time.Now()})
		}
		if v.CMDB_ID != "" {
			updateRec = append(updateRec, bson.E{Key: "CMDB_ID", Value: v.CMDB_ID})
		}
		if v.CMDB_App_Name != "" {
			updateRec = append(updateRec, bson.E{Key: "CMDB_App_Name", Value: v.CMDB_App_Name})
		}
		if v.Description != "" {
			updateRec = append(updateRec, bson.E{Key: "Description", Value: v.Description})
		}
		if v.Support_Group != "" {
			updateRec = append(updateRec, bson.E{Key: "Support_Group", Value: v.Support_Group})
		}
		if v.IT_Owner != "" {
			updateRec = append(updateRec, bson.E{Key: "IT_Owner", Value: v.IT_Owner})
			updateRec = append(updateRec, bson.E{Key: "Assigned_To", Value: v.IT_Owner})
		}
		if v.App_Business_Owner != "" {
			updateRec = append(updateRec, bson.E{Key: "App_Business_Owner", Value: v.App_Business_Owner})
		}
		if v.Technical_SME != "" {
			updateRec = append(updateRec, bson.E{Key: "Technical_SME", Value: v.Technical_SME})
		}
		if v.Criticality != "" {
			updateRec = append(updateRec, bson.E{Key: "Criticality", Value: v.Criticality})
		}
		if v.Tower_Lead != "" {
			updateRec = append(updateRec, bson.E{Key: "Tower_Lead", Value: v.Tower_Lead})
		} else {
			updateRec = append(updateRec, bson.E{Key: "Tower_Lead", Value: "Other"})
		}
		if v.Status != "" {
			updateRec = append(updateRec, bson.E{Key: "Status", Value: v.Status})
		} else {
			updateRec = append(updateRec, bson.E{Key: "Status", Value: "Not Started"})
		}
		if v.Supports_SSO != "" {
			updateRec = append(updateRec, bson.E{Key: "Supports_SSO", Value: v.Supports_SSO})
		} else {
			updateRec = append(updateRec, bson.E{Key: "Supports_SSO", Value: "No"})
		}

		if v.In_App_Auth != "" {
			updateRec = append(updateRec, bson.E{Key: "In_App_Auth", Value: v.In_App_Auth})
		} else {
			updateRec = append(updateRec, bson.E{Key: "In_App_Auth", Value: "No"})
		}

		if v.Birth_Right_Access != "" {
			updateRec = append(updateRec, bson.E{Key: "Birth_Right_Access", Value: v.Birth_Right_Access})
		} else {
			updateRec = append(updateRec, bson.E{Key: "Birth_Right_Access", Value: "No"})
		}

		if v.Birth_Right_Access_Criteria_Identified != "" {
			updateRec = append(updateRec, bson.E{Key: "Birth_Right_Access_Criteria_Identified", Value: v.Birth_Right_Access_Criteria_Identified})
		} else {
			updateRec = append(updateRec, bson.E{Key: "Birth_Right_Access_Criteria_Identified", Value: "No"})
		}

		if v.Centralized_Store_Authentication_Authorization != "" {
			updateRec = append(updateRec, bson.E{Key: "Centralized_Store_Authentication_Authorization", Value: v.Centralized_Store_Authentication_Authorization})
		} else {
			updateRec = append(updateRec, bson.E{Key: "Centralized_Store_Authentication_Authorization", Value: "No"})
		}

		if v.Multiple_Environments != "" {
			updateRec = append(updateRec, bson.E{Key: "Multiple_Environments", Value: v.Multiple_Environments})
		} else {
			updateRec = append(updateRec, bson.E{Key: "Multiple_Environments", Value: "No"})
		}

		if v.Account_Deletion_On_DeProvisioning != "" {
			updateRec = append(updateRec, bson.E{Key: "Account_Deletion_On_DeProvisioning", Value: v.Account_Deletion_On_DeProvisioning})
		} else {
			updateRec = append(updateRec, bson.E{Key: "Account_Deletion_On_DeProvisioning", Value: "No"})
		}

		if v.SOX_Status != "" {
			updateRec = append(updateRec, bson.E{Key: "SOX_Status", Value: v.SOX_Status})
		} else {
			updateRec = append(updateRec, bson.E{Key: "SOX_Status", Value: "No"})
		}
		if v.Access_Certification_Tool != "" {
			updateRec = append(updateRec, bson.E{Key: "Access_Certification_Tool", Value: v.Access_Certification_Tool})
		} else {
			updateRec = append(updateRec, bson.E{Key: "Access_Certification_Tool", Value: "Sailpoint"})
		}
		if v.Access_Request_Tool != "" {
			updateRec = append(updateRec, bson.E{Key: "Access_Request_Tool", Value: v.Access_Request_Tool})
		} else {
			updateRec = append(updateRec, bson.E{Key: "Access_Request_Tool", Value: "Sailpoint"})
		}
		if v.Access_Provisioning_Tool != "" {
			updateRec = append(updateRec, bson.E{Key: "Access_Provisioning_Tool", Value: v.Access_Provisioning_Tool})
		}
		if v.Access_Request_Approvals != "" {
			updateRec = append(updateRec, bson.E{Key: "Access_Request_Approvals", Value: v.Access_Request_Approvals})
		} else {
			updateRec = append(updateRec, bson.E{Key: "Access_Request_Approvals", Value: "Manager and App Owner"})
		}
		if v.Authentication_Tool != "" {
			updateRec = append(updateRec, bson.E{Key: "Authentication_Tool", Value: v.Authentication_Tool})
		}
		if v.MFA_Tool != "" {
			updateRec = append(updateRec, bson.E{Key: "MFA_Tool", Value: v.MFA_Tool})
		}
		if v.Current_Sailpoint_Integraton_Method != "" {
			updateRec = append(updateRec, bson.E{Key: "Current_Sailpoint_Integraton_Method", Value: v.Current_Sailpoint_Integraton_Method})
		}
		if v.Type != "" {
			updateRec = append(updateRec, bson.E{Key: "type", Value: v.Type})
		}

		setUpdate := bson.D{{Key: "$set", Value: updateRec}}

		if len(setUpdate[0].Value.(bson.D)) > 0 {
			if v.ID == primitive.NilObjectID {
				writeMod := mongo.NewUpdateOneModel().SetFilter(bson.D{{Key: "_id", Value: primitive.NewObjectID()}}).SetUpdate(setUpdate).SetUpsert(true)
				models = append(models, writeMod)
			} else {
				writeMod := mongo.NewUpdateOneModel().SetFilter(bson.D{{Key: "_id", Value: v.ID}}).SetUpdate(setUpdate).SetUpsert(true)
				models = append(models, writeMod)
			}
		}

	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := collectionAP.BulkWrite(context.Background(), models, opts)
	return res, err
}

func getAllApplicationRecordsByOwner(id string) ([]primitive.M, error) {
	var results []primitive.M
	qry := []bson.M{
		{
			// "$match": bson.M{"type": "Application",
			// 	"$or": []bson.M{
			// 		bson.M{"IT_Owner": id},
			// 		bson.M{"App_Business_Owner": id},
			// 	},
			// },
			"$match": bson.M{"type": "Application",
				"$or": []bson.M{
					bson.M{"IT_Owner": id},
					bson.M{"Technical_SME": id},
				},
			},
		},
		// {
		// 	"$lookup": bson.M{
		// 		"from":         "Application_Group",    // Child collection to join
		// 		"localField":   "application_group_id", // Parent collection reference holding field
		// 		"foreignField": "_id",                  // Child collection reference field
		// 		"as":           "Application_Group",    // Arbitrary field name to store result set
		// 	},
		// },
		{
			"$unwind": bson.M{
				"path":                       "$Application_Group",
				"preserveNullAndEmptyArrays": true,
			},
		},
		// {
		// 	"$lookup": bson.M{
		// 		"from":         "Business_Unit",                      // Child collection to join
		// 		"localField":   "Application_Group.business_unit_id", // Parent collection reference holding field
		// 		"foreignField": "_id",                                // Child collection reference field
		// 		"as":           "Business_Unit",                      // Arbitrary field name to store result set
		// 	},
		// },
		{
			"$unwind": bson.M{
				"path":                       "$Business_Unit",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	cur, err := collectionAP.Aggregate(context.Background(), qry)

	if err != nil {
		return results, err
	}
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			continue
		}
		results = append(results, result)
	}
	cur.Close(context.Background())
	return results, err
}

func GetAllApplicationByOwner(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	payload, err := getAllApplicationRecordsByOwner(params["id"])
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		json.NewEncoder(w).Encode(payload)
	}
}

// #endregion

// #region Logs
func GetAllLog(w http.ResponseWriter, r *http.Request) {
	payload, err := getAllLogRecords()
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		json.NewEncoder(w).Encode(resp)
	} else if payload == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(payload)
	}
}
func GetLog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	if id == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	payload, err := getOneLog(id)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else if payload == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(payload)
	}
}
func DeleteLog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	if id == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	deleteResult, err := deleteOneLog(id)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		json.NewEncoder(w).Encode(deleteResult)
	}
}
func getAllLogRecords() ([]primitive.M, error) {
	var results []primitive.M
	qry := []bson.M{}
	cur, err := collectionLG.Aggregate(context.Background(), qry)
	if err != nil {
		return results, err
	}
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			continue
		}
		results = append(results, result)
	}
	cur.Close(context.Background())
	return results, err
}
func getOneLog(id primitive.ObjectID) (primitive.M, error) {
	var result primitive.M
	qry := []bson.M{
		{
			"$match": bson.M{
				"_id": id,
			},
		},
	}
	cur, err := collectionLG.Aggregate(context.Background(), qry)
	if err != nil {
		return result, err
	}
	if cur.Next(context.Background()) {
		e := cur.Decode(&result)
		if e != nil {
			return result, err
		}
	}
	return result, err
}
func deleteOneLog(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{
		"$and": bson.A{bson.M{"_id": id}},
	}
	d, err := collectionLG.DeleteOne(context.Background(), filter)
	return d, err
}
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(32 << 20) // 10 MB maximum file size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Get a slice of all the uploaded files

	ids := []FileType{}
	for _, values := range r.MultipartForm.File {
		for _, file := range values {
			// Open the uploaded file
			f, err := file.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Upload the file to GridFS
			uploadStream, err := bucket.OpenUploadStream(file.Filename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer uploadStream.Close()
			if _, err := io.Copy(uploadStream, f); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Get the ID of the uploaded file
			var fT FileType
			fT.Id = uploadStream.FileID.(primitive.ObjectID).Hex()
			fT.Filename = file.Filename
			ids = append(ids, fT)
		}
	}
	// Return the IDs of the uploaded files as the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ids)
}
func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID of the file to retrieve from the request parameters
	fileID := mux.Vars(r)["id"]
	// Open a GridFS bucket
	// Retrieve the file from GridFS
	fileObjID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	downloadStream, err := bucket.OpenDownloadStream(fileObjID)
	if err != nil {
		panic(err)
	}

	fileBytes := make([]byte, 1024)
	if _, err := downloadStream.Read(fileBytes); err != nil {
		panic(err)
	}
	w.Write(fileBytes)
}

// User Operations

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

var jwtSecret = []byte("Mizuho_Secret")

func generateAccessToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"_id":         user.ID.Hex(),
		"email":       user.Email,
		"name":        user.Name,
		"companyRole": user.CompanyRole,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password = hashPassword(user.Password)

	var existingUser models.User
	err = collectionUsers.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		http.Error(w, "User with this email already exists", http.StatusBadRequest)
		return
	}

	user.UpdatedAt = time.Now()
	user.CreatedAt = time.Now()
	user.Avatar = "default_avatar_link_here"

	user.CompanyRole = []string{}

	_, err = collectionUsers.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User registered successfully"}`))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password = hashPassword(user.Password)

	err = collectionUsers.FindOne(context.TODO(), bson.M{
		"email":    user.Email,
		"password": user.Password,
	}).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := generateAccessToken(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = collectionUsers.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"accessToken": accessToken,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"accessToken": accessToken,
		"user":        user,
	})
}

func ListAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	cur, err := collectionUsers.Find(context.TODO(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		http.Error(w, "Authorization token is missing", http.StatusUnauthorized)
		return
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Error parsing claims", http.StatusInternalServerError)
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims["_id"].(string))
	if err != nil {
		http.Error(w, "Error converting ID", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:          userID,
		Email:       claims["email"].(string),
		Name:        claims["name"].(string),
		CompanyRole: convertInterfaceToStringSlice(claims["companyRole"]),
	}

	json.NewEncoder(w).Encode(user)
}

func convertInterfaceToStringSlice(val interface{}) []string {
	var result []string
	if arr, ok := val.([]interface{}); ok {
		for _, v := range arr {
			if str, ok := v.(string); ok {
				result = append(result, str)
			}
		}
	}
	return result
}

func generateRandomCode() string {
	const charset = "0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func storeResetCodeInDatabase(email, resetCode string) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"resetCode": resetCode}}
	_, err := collectionUsers.UpdateOne(context.TODO(), filter, update)
	return err
}

func getResetCodeFromDatabase(email string) (string, error) {
	var user models.User
	err := collectionUsers.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", err
	}
	return user.ResetCode, nil
}

func SendResetPasswordEmail(email, resetCode string) error {
	auth := smtp.PlainAuth("", "your_mail@gmail.com", "App_password", "smtp.gmail.com")

	from := "your_mail@gmail.com"
	to := []string{email}
	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: Password Reset Code\r\n"+
			"\r\n"+
			"Here is your password reset code: %s\r\n",
		email, resetCode))

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		log.Printf("Error sending email : %s", err)
		return err
	}

	log.Printf("Email sent to %s", email)
	return nil
}

func RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := user.Email

	var existingUser models.User
	err = collectionUsers.FindOne(context.TODO(), bson.M{"email": email}).Decode(&existingUser)
	if err != nil {
		http.Error(w, "User with this email does not exist", http.StatusBadRequest)
		return
	}
	resetCode := generateRandomCode()
	err = storeResetCodeInDatabase(email, resetCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = SendResetPasswordEmail(email, resetCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Reset code sent successfully",
	})
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	email := user.Email
	resetCode := user.ResetCode
	newPassword := user.Password

	storedResetCode, err := getResetCodeFromDatabase(email)
	if err != nil || storedResetCode != resetCode {
		http.Error(w, "Invalid reset code/email", http.StatusBadRequest)
		return
	}

	newHashedPassword := hashPassword(newPassword)
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"password": newHashedPassword, "resetCode": ""}}
	_, err = collectionUsers.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Password reset successful",
	})
}

func ResetPasswordForUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	email := user.Email
	newPassword := user.Password

	err = collectionUsers.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid Email", http.StatusBadRequest)
	}

	newHashedPassword := hashPassword(newPassword)
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"password": newHashedPassword}}
	_, err = collectionUsers.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Password reset successful",
	})
}

func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": id}

	if user.Name != "" && len(user.CompanyRole) != 0 {
		update := bson.M{
			"$set": bson.M{
				"companyRole": user.CompanyRole,
				"name":        user.Name,
			},
		}
		_, err := collectionUsers.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {

		http.Error(w, "Couldn't update user details!!", http.StatusBadRequest)
		return

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User role update successfully!!!",
	})

}

func deleteOneUser(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	d, err := collectionUsers.DeleteOne(context.Background(), filter)
	return d, err
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	if id == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		var resp ErrorResult
		resp.Error.Message = "Missing ID"
		json.NewEncoder(w).Encode(resp)
		return
	}
	_, err := deleteOneUser(id)
	if err != nil {
		var resp ErrorResult
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User deleted successfully!!!",
		})
	}
}

func UpdateCertification(Certification models.Certification) (*mongo.InsertOneResult, error) {
	insertResult, err := collectionCert.InsertOne(context.Background(), Certification)
	return insertResult, err
}
