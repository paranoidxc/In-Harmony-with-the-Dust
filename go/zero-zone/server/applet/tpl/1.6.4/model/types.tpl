type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
		FindAllByWhere(ctx context.Context, where string) ([]*{{.upperStartCamelObject}}, error)
        FindAllByWhereCount(ctx context.Context, where string) (int64, error)
	    FindPageByWhere(ctx context.Context, where string, page int64, limit int64) ([]*{{.upperStartCamelObject}}, error)
		FindPageByWhereCount(ctx context.Context, where string) (int64, error)	}

	default{{.upperStartCamelObject}}Model struct {
		{{if .withCache}}sqlc.CachedConn{{else}}conn sqlx.SqlConn{{end}}
		table string
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
