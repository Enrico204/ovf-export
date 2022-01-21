package manifest

import "testing"

func TestParse_OK_vbox(t *testing.T) {
	var fcontent = `SHA1 (Mikrotik.ovf) = 8cd5371c4042ab1a11f9963d8eba43c99b8d1f8e
SHA1 (Mikrotik-disk001.vmdk) = 9bd8182e143b977e96d1c9241a1fbcdc9252aa72`

	mf, err := Parse([]byte(fcontent))
	if err != nil {
		t.Fatal(err)
	}
	switch {
	case len(mf) != 2:
		t.Fatal("Expected 2 items, found ", len(mf))
	case mf[0].Algorithm != "SHA1":
		t.Fatal("Invalid algorithm for 1st item, expected SHA1, found ", mf[0].Algorithm)
	case mf[0].Name != "Mikrotik.ovf":
		t.Fatal("Invalid name for 1st item, expected Mikrotik.ovf, found ", mf[0].Name)
	case mf[0].Hash != "8cd5371c4042ab1a11f9963d8eba43c99b8d1f8e":
		t.Fatal("Invalid hash for 1st item, expected 8cd5371c4042ab1a11f9963d8eba43c99b8d1f8e, found ", mf[0].Hash)
	case mf[1].Algorithm != "SHA1":
		t.Fatal("Invalid algorithm for 2nd item, expected SHA1, found ", mf[1].Algorithm)
	case mf[1].Name != "Mikrotik-disk001.vmdk":
		t.Fatal("Invalid name for 2nd item, expected Mikrotik-disk001.vmdk, found ", mf[1].Name)
	case mf[1].Hash != "9bd8182e143b977e96d1c9241a1fbcdc9252aa72":
		t.Fatal("Invalid hash for 2nd item, expected 9bd8182e143b977e96d1c9241a1fbcdc9252aa72, found ", mf[1].Hash)
	}
}

func TestParse_OK_vmware(t *testing.T) {
	var fcontent = `SHA1(Mikrotik.ovf)= 8cd5371c4042ab1a11f9963d8eba43c99b8d1f8e
SHA1(Mikrotik-disk001.vmdk)= 9bd8182e143b977e96d1c9241a1fbcdc9252aa72`

	mf, err := Parse([]byte(fcontent))
	if err != nil {
		t.Fatal(err)
	}
	switch {
	case len(mf) != 2:
		t.Fatal("Expected 2 items, found ", len(mf))
	case mf[0].Algorithm != "SHA1":
		t.Fatal("Invalid algorithm for 1st item, expected SHA1, found ", mf[0].Algorithm)
	case mf[0].Name != "Mikrotik.ovf":
		t.Fatal("Invalid name for 1st item, expected Mikrotik.ovf, found ", mf[0].Name)
	case mf[0].Hash != "8cd5371c4042ab1a11f9963d8eba43c99b8d1f8e":
		t.Fatal("Invalid hash for 1st item, expected 8cd5371c4042ab1a11f9963d8eba43c99b8d1f8e, found ", mf[0].Hash)
	case mf[1].Algorithm != "SHA1":
		t.Fatal("Invalid algorithm for 2nd item, expected SHA1, found ", mf[1].Algorithm)
	case mf[1].Name != "Mikrotik-disk001.vmdk":
		t.Fatal("Invalid name for 2nd item, expected Mikrotik-disk001.vmdk, found ", mf[1].Name)
	case mf[1].Hash != "9bd8182e143b977e96d1c9241a1fbcdc9252aa72":
		t.Fatal("Invalid hash for 2nd item, expected 9bd8182e143b977e96d1c9241a1fbcdc9252aa72, found ", mf[1].Hash)
	}
}

func TestParse_badFormat(t *testing.T) {
	var fcontent = `Mikrotik.ovf 8cd5371c4042ab1a11f9963d8eba43c99b8d1f8e`

	_, err := Parse([]byte(fcontent))
	if err == nil {
		t.Fatal("expected parse error, found ok")
	}
}

func TestParse_badFormat2(t *testing.T) {
	var fcontent = `Mikrotik.ovf = 8cd5371c4042ab1a11f9963d8eba43c99b8d1f8e`

	_, err := Parse([]byte(fcontent))
	if err == nil {
		t.Fatal("expected parse error, found ok")
	}
}

func TestParse_badFormat3(t *testing.T) {
	var fcontent = `Mikro(tik).ovf = 8cd5371c4042ab1a11f9963d8eba43c99b8d1f8e`

	_, err := Parse([]byte(fcontent))
	if err == nil {
		t.Fatal("expected parse error, found ok")
	}
}

func TestParse_badFormat4(t *testing.T) {
	var fcontent = `8cd5371c4042ab1a11f9963d8eba43c99b8=asdasdd1f8e`

	_, err := Parse([]byte(fcontent))
	if err == nil {
		t.Fatal("expected parse error, found ok")
	}
}
