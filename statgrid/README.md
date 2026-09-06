# statgrid

`statgrid` is a responsive, presentation-only grid component for displaying key statistics/metrics in key-value pairs (e.g. "80+" / "Años de historia").

## Usage

```go
import "webtyp.com/components/statgrid"

grid := &statgrid.StatGrid{
    Items: []statgrid.StatItem{
        {Value: "80+", Label: "Años de historia"},
        {Value: "15,000+", Label: "Pacientes atendidos"},
        {Value: "24/7", Label: "Atención continuada"},
    },
}
```
