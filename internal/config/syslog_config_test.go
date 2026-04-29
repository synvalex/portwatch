package config

import "testing"

func TestDefaultSyslogConfig(t *testing.T) {
	c := DefaultSyslogConfig()
	if c.Enabled {
		t.Error("expected Enabled=false")
	}
	if c.Tag != "portwatch" {
		t.Errorf("expected tag=portwatch, got %q", c.Tag)
	}
	if c.Facility != "daemon" {
		t.Errorf("expected facility=daemon, got %q", c.Facility)
	}
}

func TestSyslogConfig_Validate_Disabled(t *testing.T) {
	c := SyslogConfig{Enabled: false}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSyslogConfig_Validate_EnabledValid(t *testing.T) {
	c := SyslogConfig{Enabled: true, Tag: "portwatch", Facility: "local0"}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSyslogConfig_Validate_EmptyTag(t *testing.T) {
	c := SyslogConfig{Enabled: true, Tag: "", Facility: "daemon"}
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty tag")
	}
}

func TestSyslogConfig_Validate_InvalidFacility(t *testing.T) {
	c := SyslogConfig{Enabled: true, Tag: "portwatch", Facility: "bogus"}
	if err := c.Validate(); err == nil {
		t.Error("expected error for invalid facility")
	}
}

func TestSyslogConfig_Validate_NetworkWithoutAddr(t *testing.T) {
	c := SyslogConfig{Enabled: true, Tag: "portwatch", Facility: "daemon", Network: "tcp", Addr: ""}
	if err := c.Validate(); err == nil {
		t.Error("expected error for network without addr")
	}
}

func TestSyslogConfig_Merge_FillsTag(t *testing.T) {
	def := DefaultSyslogConfig()
	c := SyslogConfig{Enabled: true, Facility: "local1"}
	merged := c.Merge(def)
	if merged.Tag != "portwatch" {
		t.Errorf("expected merged tag=portwatch, got %q", merged.Tag)
	}
}

func TestSyslogConfig_Merge_PreservesUserTag(t *testing.T) {
	def := DefaultSyslogConfig()
	c := SyslogConfig{Enabled: true, Tag: "myapp", Facility: "daemon"}
	merged := c.Merge(def)
	if merged.Tag != "myapp" {
		t.Errorf("expected merged tag=myapp, got %q", merged.Tag)
	}
}
