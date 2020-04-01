package sfAuth

type EtcdHelpter struct {
    EtcdUrl    string
    EtcdCaPath string
}

func (e *EtcdHelpter) Get(key string) (string, error) {
    return "", nil
}

func (e *EtcdHelpter) AddOrUpdate(key string, value string) (string, error) {
    return "", nil
}