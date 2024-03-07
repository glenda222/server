package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	Name        string             `json:"name" bson:"name"`
	CompanyRole []string           `json:"companyRole" bson:"companyRole"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	Avatar      string             `json:"avatar" bson:"avatar"`
	AccessToken string             `json:"accessToken" bson:"accessToken"`
	ResetCode   string             `json:"resetCode" bson:"resetCode"`
}

type LogType struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Log     string             `json:"log,omitempty" bson:"log,omitempty"`
	Created time.Time          `json:"created,omitempty" bson:"created,omitempty"`
}

type BusinessUnitType struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type        string             `json:"type,omitempty" bson:"type,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Updated     time.Time          `json:"updated,omitempty" bson:"updated,omitempty"`
	Modified_By string             `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
}

type ApplicationGroupType struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type             string             `json:"type,omitempty" bson:"type,omitempty"`
	Name             string             `json:"name,omitempty" bson:"name,omitempty"`
	Updated          time.Time          `json:"updated,omitempty" bson:"updated,omitempty"`
	Modified_By      string             `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	Business_Unit_ID primitive.ObjectID `json:"business_unit_id,omitempty" bson:"business_unit_id,omitempty"`
}
type Migration_Plan struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	X  string             `json:"x,omitempty" bson:"x,omitempty"`
	Y  []string           `json:"y,omitempty" bson:"y,omitempty"`
}

type Attribute_Mapping struct {
	ID                              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Application_Attribute_Name      string             `json:"Application_Attribute_Name,omitempty" bson:"Application_Attribute_Name,omitempty"`
	Display_Name                    string             `json:"Display_Name,omitempty" bson:"Display_Name,omitempty"`
	Mandatory_Attribute             string             `json:"Mandatory_Attribute,omitempty" bson:"Mandatory_Attribute,omitempty"`
	Data_Type                       string             `json:"Data_Type,omitempty" bson:"Data_Type,omitempty"`
	Attribute_Rule_Description      string             `json:"Attribute_Rule_Description,omitempty" bson:"Attribute_Rule_Description,omitempty"`
	Comma_Separated_Possible_Values string             `json:"Comma_Separated_Possible_Values,omitempty" bson:"Comma_Separated_Possible_Values,omitempty"`
}
type SOD_Details struct {
	App_1         string `json:"App_1,omitempty" bson:"App_1,omitempty"`
	Entitlement_1 string `json:"Entitlement_1,omitempty" bson:"Entitlement_1,omitempty"`
	App_2         string `json:"App_2,omitempty" bson:"App_2,omitempty"`
	Entitlement_2 string `json:"Entitlement_2,omitempty" bson:"Entitlement_2,omitempty"`
}

type ApplicationType struct {
	ID                                             primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Type                                           string               `json:"type,omitempty" bson:"type,omitempty"`
	Name                                           string               `json:"name,omitempty" bson:"name,omitempty"`
	Updated                                        time.Time            `json:"updated,omitempty" bson:"updated,omitempty"`
	Modified_By                                    string               `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	Application_Group_ID                           primitive.ObjectID   `json:"application_group_id,omitempty" bson:"application_group_id,omitempty"`
	Assigned_To                                    string               `json:"Assigned_To,omitempty" bson:"Assigned_To,omitempty"`
	Migration_Plan                                 []Migration_Plan     `json:"Migration_Plan,omitempty" bson:"Migration_Plan,omitempty"`
	Attribute_Mapping                              []Attribute_Mapping  `json:"Attribute_Mapping,omitempty" bson:"Attribute_Mapping,omitempty"`
	SOD_Details                                    []SOD_Details        `json:"SOD_Details,omitempty" bson:"SOD_Details,omitempty"`
	Attached_Files                                 []FileType           `json:"Attached_Files,omitempty" bson:"Attached_Files,omitempty"`
	CMDB_ID                                        string               `json:"CMDB_ID,omitempty" bson:"CMDB_ID,omitempty"`
	Description                                    string               `json:"Description,omitempty" bson:"Description,omitempty"`
	Support_Group                                  string               `json:"Support_Group,omitempty" bson:"Support_Group,omitempty"`
	App_Name                                       string               `json:"App_Name,omitempty" bson:"App_Name,omitempty"`
	CMDB_App_Name                                  string               `json:"CMDB_App_Name,omitempty" bson:"CMDB_App_Name,omitempty"`
	IT_Owner                                       string               `json:"IT_Owner,omitempty" bson:"IT_Owner,omitempty"`
	App_Business_Owner                             string               `json:"App_Business_Owner,omitempty" bson:"App_Business_Owner,omitempty"`
	Technical_SME                                  string               `json:"Technical_SME,omitempty" bson:"Technical_SME,omitempty"`
	Tower_Lead                                     string               `json:"Tower_Lead,omitempty" bson:"Tower_Lead,omitempty"`
	Status                                         string               `json:"Status,omitempty" bson:"Status,omitempty"`
	SOX_Status                                     string               `json:"SOX_Status,omitempty" bson:"SOX_Status,omitempty"`
	Criticality                                    string               `json:"Criticality,omitempty" bson:"Criticality,omitempty"`
	Access_Certification_Tool                      string               `json:"Access_Certification_Tool,omitempty" bson:"Access_Certification_Tool,omitempty"`
	Access_Request_Tool                            string               `json:"Access_Request_Tool,omitempty" bson:"Access_Request_Tool,omitempty"`
	Access_Provisioning_Tool                       string               `json:"Access_Provisioning_Tool,omitempty" bson:"Access_Provisioning_Tool,omitempty"`
	Authentication_Tool                            string               `json:"Authentication_Tool,omitempty" bson:"Authentication_Tool,omitempty"`
	MFA_Tool                                       string               `json:"MFA_Tool,omitempty" bson:"MFA_Tool,omitempty"`
	Application_Type                               string               `json:"Application_Type,omitempty" bson:"Application_Type,omitempty"`
	Application_Deployment_Style                   string               `json:"Application_Deployment_Style,omitempty" bson:"Application_Deployment_Style,omitempty"`
	Centralized_Store_Authentication_Authorization string               `json:"Centralized_Store_Authentication_Authorization,omitempty" bson:"Centralized_Store_Authentication_Authorization,omitempty"`
	Directory_Services_Used                        string               `json:"Directory_Services_Used,omitempty" bson:"Directory_Services_Used,omitempty"`
	Birth_Right_Access                             string               `json:"Birth_Right_Access,omitempty" bson:"Birth_Right_Access,omitempty"`
	Birth_Right_Access_Criteria_Identified         string               `json:"Birth_Right_Access_Criteria_Identified,omitempty" bson:"Birth_Right_Access_Criteria_Identified,omitempty"`
	Provisioning_Pre_requisites                    string               `json:"Provisioning_Pre_requisites,omitempty" bson:"Provisioning_Pre_requisites,omitempty"`
	Post_Termination_Process_Requirements          string               `json:"Post_Termination_Process_Requirements,omitempty" bson:"Post_Termination_Process_Requirements,omitempty"`
	Supports_SSO                                   string               `json:"Supports_SSO,omitempty" bson:"Supports_SSO,omitempty"`
	Open_Standards_Support                         []string             `json:"Open_Standards_Support,omitempty" bson:"Open_Standards_Support,omitempty"`
	In_App_Auth                                    string               `json:"In_App_Auth,omitempty" bson:"In_App_Auth,omitempty"`
	Multiple_Environments                          string               `json:"Multiple_Environments,omitempty" bson:"Multiple_Environments,omitempty"`
	Application_User_Domain                        string               `json:"Application_User_Domain,omitempty" bson:"Application_User_Domain,omitempty"`
	Access_Request_Approvals                       string               `json:"Access_Request_Approvals,omitempty" bson:"Access_Request_Approvals,omitempty"`
	Other_Access_Request_Approvals                 string               `json:"Other_Access_Request_Approvals,omitempty" bson:"Other_Access_Request_Approvals,omitempty"`
	User_Population_For_App                        string               `json:"User_Population_For_App,omitempty" bson:"User_Population_For_App,omitempty"`
	Account_Deletion_On_DeProvisioning             string               `json:"Account_Deletion_On_DeProvisioning,omitempty" bson:"Account_Deletion_On_DeProvisioning,omitempty"`
	Process_To_DeProvision                         string               `json:"Process_To_DeProvision,omitempty" bson:"Process_To_DeProvision,omitempty"`
	Unique_Identifier_Generation_Logic             string               `json:"Unique_Identifier_Generation_Logic,omitempty" bson:"Unique_Identifier_Generation_Logic,omitempty"`
	Unique_Identifier_For_User                     string               `json:"Unique_Identifier_For_User,omitempty" bson:"Unique_Identifier_For_User,omitempty"`
	Status_Attribute_Name                          string               `json:"Status_Attribute_Name,omitempty" bson:"Status_Attribute_Name,omitempty"`
	DEV_Endpoint_Base_URL                          string               `json:"DEV_Endpoint_Base_URL,omitempty" bson:"DEV_Endpoint_Base_URL,omitempty"`
	UAT_Endpoint_Base_URL                          string               `json:"UAT_Endpoint_Base_URL,omitempty" bson:"UAT_Endpoint_Base_URL,omitempty"`
	PROD_Endpoint_Base_URL                         string               `json:"PROD_Endpoint_Base_URL,omitempty" bson:"PROD_Endpoint_Base_URL,omitempty"`
	DEV_Auth_Type                                  string               `json:"DEV_Auth_Type,omitempty" bson:"DEV_Auth_Type,omitempty"`
	UAT_Auth_Type                                  string               `json:"UAT_Auth_Type,omitempty" bson:"UAT_Auth_Type,omitempty"`
	PROD_Auth_Type                                 string               `json:"PROD_Auth_Type,omitempty" bson:"PROD_Auth_Type,omitempty"`
	DEV_Client_ID                                  string               `json:"DEV_Client_ID,omitempty" bson:"DEV_Client_ID,omitempty"`
	UAT_Client_ID                                  string               `json:"UAT_Client_ID,omitempty" bson:"UAT_Client_ID,omitempty"`
	PROD_Client_ID                                 string               `json:"PROD_Client_ID,omitempty" bson:"PROD_Client_ID,omitempty"`
	DEV_Client_Secret                              string               `json:"DEV_Client_Secret,omitempty" bson:"DEV_Client_Secret,omitempty"`
	UAT_Client_Secret                              string               `json:"UAT_Client_Secret,omitempty" bson:"UAT_Client_Secret,omitempty"`
	PROD_Client_Secret                             string               `json:"PROD_Client_Secret,omitempty" bson:"PROD_Client_Secret,omitempty"`
	DEV_JDBC_URL                                   string               `json:"DEV_JDBC_URL,omitempty" bson:"DEV_JDBC_URL,omitempty"`
	UAT_JDBC_URL                                   string               `json:"UAT_JDBC_URL,omitempty" bson:"UAT_JDBC_URL,omitempty"`
	PROD_JDBC_URL                                  string               `json:"PROD_JDBC_URL,omitempty" bson:"PROD_JDBC_URL,omitempty"`
	DEV_JDBC_Driver                                string               `json:"DEV_JDBC_Driver,omitempty" bson:"DEV_JDBC_Driver,omitempty"`
	UAT_JDBC_Driver                                string               `json:"UAT_JDBC_Driver,omitempty" bson:"UAT_JDBC_Driver,omitempty"`
	PROD_JDBC_Driver                               string               `json:"PROD_JDBC_Driver,omitempty" bson:"PROD_JDBC_Driver,omitempty"`
	DEV_LDAP_URL                                   string               `json:"DEV_LDAP_URL,omitempty" bson:"DEV_LDAP_URL,omitempty"`
	UAT_LDAP_URL                                   string               `json:"UAT_LDAP_URL,omitempty" bson:"UAT_LDAP_URL,omitempty"`
	PROD_LDAP_URL                                  string               `json:"PROD_LDAP_URL,omitempty" bson:"PROD_LDAP_URL,omitempty"`
	DEV_File_Path                                  string               `json:"DEV_File_Path,omitempty" bson:"DEV_File_Path,omitempty"`
	UAT_File_Path                                  string               `json:"UAT_File_Path,omitempty" bson:"UAT_File_Path,omitempty"`
	PROD_File_Path                                 string               `json:"PROD_File_Path,omitempty" bson:"PROD_File_Path,omitempty"`
	DEV_Delimiter                                  string               `json:"DEV_Delimiter,omitempty" bson:"DEV_Delimiter,omitempty"`
	UAT_Delimiter                                  string               `json:"UAT_Delimiter,omitempty" bson:"UAT_Delimiter,omitempty"`
	PROD_Delimiter                                 string               `json:"PROD_Delimiter,omitempty" bson:"PROD_Delimiter,omitempty"`
	DEV_Username                                   string               `json:"DEV_Username,omitempty" bson:"DEV_Username,omitempty"`
	UAT_Username                                   string               `json:"UAT_Username,omitempty" bson:"UAT_Username,omitempty"`
	PROD_Username                                  string               `json:"PROD_Username,omitempty" bson:"PROD_Username,omitempty"`
	DEV_Password                                   string               `json:"DEV_Password,omitempty" bson:"DEV_Password,omitempty"`
	UAT_Password                                   string               `json:"UAT_Password,omitempty" bson:"UAT_Password,omitempty"`
	PROD_Password                                  string               `json:"PROD_Password,omitempty" bson:"PROD_Password,omitempty"`
	Connector_Integraton_Framework                 string               `json:"Connector_Integraton_Framework,omitempty" bson:"Connector_Integraton_Framework,omitempty"`
	Other_Integraton_Method                        string               `json:"Other_Integraton_Method,omitempty" bson:"Other_Integraton_Method,omitempty"`
	JDBC_Env_Details                               string               `json:"JDBC_Env_Details" bson:"JDBC_Env_Details,omitempty"`
	JDBC_ServiceAccountName                        string               `json:"JDBC_ServiceAccountName" bson:"JDBC_ServiceAccountName,omitempty"`
	JDBC_Host                                      string               `json:"JDBC_Host" bson:"JDBC_Host,omitempty"`
	JDBC_Port                                      string               `json:"JDBC_Port" bson:"JDBC_Port,omitempty"`
	Query_Import_UserAccounts                      string               `json:"Query_Import_UserAccounts" bson:"Query_Import_UserAccounts,omitempty"`
	Query_Create_NewAccount                        string               `json:"Query_Create_NewAccount" bson:"Query_Create_NewAccount,omitempty"`
	Query_Import_Roles                             string               `json:"Query_Import_Roles" bson:"Query_Import_Roles,omitempty"`
	Query_Assign_Role                              string               `json:"Query_Assign_Role" bson:"Query_Assign_Role,omitempty"`
	Query_Remove_Role                              string               `json:"Query_Remove_Role" bson:"Query_Remove_Role,omitempty"`
	Query_Disable_Account                          string               `json:"Query_Disable_Account" bson:"Query_Disable_Account,omitempty"`
	Query_Delete_Account                           string               `json:"Query_Delete_Account" bson:"Query_Delete_Account,omitempty"`
	Current_Sailpoint_Integraton_Method            string               `json:"Current_Sailpoint_Integraton_Method" bson:"Current_Sailpoint_Integraton_Method"`
	Certified                                      bool                 `json:"certified" bson:"certified,omitempty"`
	LogIds                                         []primitive.ObjectID `json:"logids" bson:"logids"`

	Api_Env_Details         string `json:"Api_Env_Details" bson:"Api_Env_Details,omitempty"`
	Api_ServiceAccountName  string `json:"Api_ServiceAccountName" bson:"Api_ServiceAccountName,omitempty"`
	Api_Host                string `json:"Api_Host" bson:"Api_Host,omitempty"`
	Api_Port                string `json:"Api_Port" bson:"Api_Port,omitempty"`
	Api_Url                 string `json:"Api_Url" bson:"Api_Url,omitempty"`
	Api_Import_UserAccounts string `json:"Api_Import_UserAccounts" bson:"Api_Import_UserAccounts,omitempty"`
	Api_Create_NewAccount   string `json:"Api_Create_NewAccount" bson:"Api_Create_NewAccount,omitempty"`
	Api_Import_Roles        string `json:"Api_Import_Roles" bson:"Api_Import_Roles,omitempty"`
	Api_Assign_Role         string `json:"Api_Assign_Role" bson:"Api_Assign_Role,omitempty"`
	Api_Remove_Role         string `json:"Api_Remove_Role" bson:"Api_Remove_Role,omitempty"`
	Api_Disable_Account     string `json:"Api_Disable_Account" bson:"Api_Disable_Account,omitempty"`
	Api_Delete_Account      string `json:"Api_Delete_Account" bson:"Api_Delete_Account,omitempty"`

	RPA_Application_Url string `json:"RPA_Application_Url" bson:"RPA_Application_Url,omitempty"`
	RPA_UserName        string `json:"RPA_UserName" bson:"RPA_UserName,omitempty"`

	LDAP_Directory_Server               string `json:"LDAP_Directory_Server" bson:"LDAP_Directory_Server,omitempty"`
	LDAP_Port                           string `json:"LDAP_Port" bson:"LDAP_Port,omitempty"`
	LDAP_Integration_ServiceAccountName string `json:"LDAP_Integration_ServiceAccountName" bson:"LDAP_Integration_ServiceAccountName,omitempty"`
	LDAP_Domain_Name                    string `json:"LDAP_Domain_Name" bson:"LDAP_Domain_Name,omitempty"`
	LDAP_Account_Group_Search_Scope     string `json:"LDAP_Account_Group_Search_Scope" bson:"LDAP_Account_Group_Search_Scope,omitempty"`

	AD_Domain_Name    string `json:"AD_Domain_Name" bson:"AD_Domain_Name,omitempty"`
	AD_Groups         string `json:"AD_Groups" bson:"AD_Groups,omitempty"`
	Descoped_Comments string `json:"Descoped_Comments,omitempty" bson:"Descoped_Comments,omitempty"`
	Additional_Notes  string `json:"Additional_Notes,omitempty" bson:"Additional_Notes,omitempty"`
}

type FileType struct {
	Id       string
	Filename string
}

type Parameter struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CloudFitmentParams []CloudFitmentParam
}

type CloudFitmentParam struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Depends    string             `json:"depends,omitempty" bson:"depends,omitempty"`
	DependsVal string             `json:"dependsVal,omitempty" bson:"dependsVal,omitempty"`
	EmptyVal   int                `json:"emptyVal,omitempty" bson:"emptyVal,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Weight     int                `json:"weight,omitempty" bson:"weight,omitempty"`
	Value      int                `json:"value,omitempty" bson:"value,omitempty"`
	Fields     []CloudFitmentParam
	Refs       map[string]interface{} `json:"refs"  bson:"refs,omitempty"`
}

type Certification struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Certified_by string             `json:"certified_by" bson:"certified_by,omitempty"`
	CMDB_ID      string             `json:"CMDB_ID" bson:"CMDB_ID,omitempty"`
	Updated      time.Time          `json:"updated,omitempty" bson:"updated,omitempty"`
}
