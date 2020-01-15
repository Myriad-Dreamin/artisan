package artisan

// VirtualService has father and sons, and this is just an abstract service that
// for deriving other class
// Getter of Virtual Service
// Get/Set
//    Faz
//    Models
//    Name
//    FilePath
type VirtualService struct {
	base     string
	models   []*model
	name     string
	filePath string
}

// Getter of Virtual Service

func (v *VirtualService) GetBase() string {
	return v.base
}

func (v *VirtualService) GetModels() []*model {
	return v.models
}

func (v *VirtualService) GetName() string {
	return v.name
}

func (v *VirtualService) GetFilePath() string {
	return v.filePath
}

// Setter of Virtual Service

func (v *VirtualService) Base(base string) ProposingService {
	v.base = base
	return v
}

func (v *VirtualService) UseModel(models ...*model) ProposingService {
	v.models = append(v.models, models...)
	return v
}

func (v *VirtualService) Name(name string) ProposingService {
	v.name = name
	return v
}

func (v *VirtualService) ToFile(fileName string) ProposingService {
	v.filePath = fileName
	return v
}