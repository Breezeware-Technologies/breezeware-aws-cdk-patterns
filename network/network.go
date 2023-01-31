// Package network provides types for handling network related resources.
//
// Contains a type for looking-up an VPC.
package network

// A VpcProps represents properties for looking-up an VPC.
//
// Configure only IsDefault field and omit others, if the VPC looking-up is a default one.
type VpcProps struct {
	Id        string `field:"optional"` // Id of the VPC
	IsDefault bool   `field:"optional"` // IsDefault flag represents whether the VPC is a default one or not. Omit other fields if default.
}
