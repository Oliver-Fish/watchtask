package execmanager

import "testing"

func TestSpawnAndKill(t *testing.T) {
	tt := map[string]*Task{
		"Non-Blocking Command (No Args)":   {CmdPath: "ls"},
		"Non-Blocking Command (With Args)": {CmdPath: "ls", Args: []string{"-lah", "~"}},
		"Blocking Command":                 {CmdPath: "ping", Args: []string{"localhost"}},
	}
	for k, v := range tt {
		t.Run(k, func(t *testing.T) {
			err := v.Start()
			if err != nil {
				t.Fatal(err)
			}
			err = v.Kill()
			if err != nil {
				t.Fatal(err)
			}
		})
	}

}
