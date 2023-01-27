package network

// VpcProps properties
type VpcProps struct {
	// Id of the VPC to lookup
	Id string `field:"optional"`
	//	// Name of the VPC to lookup
	//	Name string `field:"optional"`
	// Is the default VPC. Omit other fields if the VPC is deafult
	IsDefault bool `field:"optional"`
}
