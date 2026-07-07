package devspace

import (
	"fmt"
	"io"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

type WorkspaceOverview struct {
	WorkspaceRoot   string                `json:"workspaceRoot"`
	ManifestVersion int                   `json:"manifestVersion"`
	ThisMachine     string                `json:"thisMachine"`
	Machines        []Machine             `json:"machines"`
	Users           []User                `json:"users,omitempty"`
	Teams           []Team                `json:"teams,omitempty"`
	Sync            WorkspaceOverviewSync `json:"sync"`
	Summary         WorkspaceStatusReport `json:"summary"`
}

type WorkspaceOverviewSync struct {
	ManifestRemote string `json:"manifestRemote,omitempty"`
	HostedEndpoint string `json:"hostedEndpoint,omitempty"`
	LastSyncAt     string `json:"lastSyncAt,omitempty"`
	LastScanAt     string `json:"lastScanAt,omitempty"`
}

func buildWorkspaceOverview() (WorkspaceOverview, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return WorkspaceOverview{}, err
	}
	m, err := LoadManifest(cfg.WorkspaceRoot)
	if err != nil {
		return WorkspaceOverview{}, err
	}
	st, err := LoadState()
	if err != nil && !missing(err) {
		return WorkspaceOverview{}, err
	}
	summary, err := buildWorkspaceStatusReport()
	if err != nil {
		return WorkspaceOverview{}, err
	}
	return WorkspaceOverview{
		WorkspaceRoot:   cfg.WorkspaceRoot,
		ManifestVersion: m.Version,
		ThisMachine:     machineLabel(cfg.MachineName, cfg.MachineID),
		Machines:        m.Machines,
		Users:           m.Users,
		Teams:           m.Teams,
		Sync: WorkspaceOverviewSync{
			ManifestRemote: redactRemote(cfg.ManifestRemote),
			HostedEndpoint: cfg.HostedSyncEndpoint,
			LastSyncAt:     st.LastSyncAt,
			LastScanAt:     st.LastScanAt,
		},
		Summary: summary,
	}, nil
}

func printWorkspaceOverview(out io.Writer, ov WorkspaceOverview) {
	out = styledWriter(out)
	fmt.Fprintln(out, currentTheme.Header.Render("Workspace"))
	fmt.Fprintf(out, "Root: %s\n", valueOrDash(ov.WorkspaceRoot))
	fmt.Fprintf(out, "Manifest version: %d\n", ov.ManifestVersion)
	fmt.Fprintf(out, "This machine: %s\n\n", valueOrDash(ov.ThisMachine))

	fmt.Fprintln(out, currentTheme.Header.Render("Machines"))
	if len(ov.Machines) == 0 {
		fmt.Fprintln(out, "No machines recorded.")
	} else {
		rows := make([][]string, 0, len(ov.Machines))
		for _, m := range ov.Machines {
			rows = append(rows, []string{valueOrDash(m.Name), valueOrDash(m.ID), valueOrDash(m.LastSeenAt)})
		}
		tbl := table.New().
			Headers("NAME", "ID", "LAST SEEN").
			Rows(rows...).
			BorderStyle(currentTheme.Muted).
			StyleFunc(func(row, _ int) lipgloss.Style {
				if row == table.HeaderRow {
					return currentTheme.Header.Padding(0, 1)
				}
				return lipgloss.NewStyle().Padding(0, 1)
			})
		fmt.Fprintln(out, tbl.Render())
	}

	fmt.Fprintf(out, "\nUsers: %s\n", namesOrDash(ov.Users))
	fmt.Fprintf(out, "Teams: %s\n\n", teamNamesOrDash(ov.Teams))

	fmt.Fprintln(out, currentTheme.Header.Render("Sync"))
	fmt.Fprintf(out, "Manifest remote: %s\n", valueOrDash(ov.Sync.ManifestRemote))
	fmt.Fprintf(out, "Hosted endpoint: %s\n", valueOrDash(ov.Sync.HostedEndpoint))
	fmt.Fprintf(out, "Last sync: %s\n", valueOrDash(ov.Sync.LastSyncAt))
	fmt.Fprintf(out, "Last scan: %s\n\n", valueOrDash(ov.Sync.LastScanAt))

	fmt.Fprintln(out, currentTheme.Header.Render("Summary"))
	fmt.Fprintf(out, "Projects tracked: %d\n", ov.Summary.ProjectsTracked)
	fmt.Fprintf(out, "Hydrated: %d\n", ov.Summary.Hydrated)
	fmt.Fprintf(out, "Placeholders: %d\n", ov.Summary.Placeholders)
	fmt.Fprintf(out, "Dirty repos: %d\n", ov.Summary.Dirty)
	fmt.Fprintf(out, "Missing env files: %d\n", ov.Summary.MissingEnv)
	fmt.Fprintf(out, "Outdated repos: %d\n", ov.Summary.Outdated)
}

func machineLabel(name, id string) string {
	name = strings.TrimSpace(name)
	id = strings.TrimSpace(id)
	if name == "" {
		return id
	}
	if id == "" {
		return name
	}
	return fmt.Sprintf("%s (%s)", name, id)
}

func namesOrDash(users []User) string {
	names := make([]string, 0, len(users))
	for _, u := range users {
		names = append(names, valueOrDash(u.Name))
	}
	if len(names) == 0 {
		return "-"
	}
	return strings.Join(names, ", ")
}

func teamNamesOrDash(teams []Team) string {
	names := make([]string, 0, len(teams))
	for _, t := range teams {
		names = append(names, valueOrDash(t.Name))
	}
	if len(names) == 0 {
		return "-"
	}
	return strings.Join(names, ", ")
}
