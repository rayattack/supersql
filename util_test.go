package supersql_test

import "github.com/rayattack/supersql"

func generateJoinCommand() supersql.Command {
	cols := []interface{}{"c.email", "c.first_name || ', ' || c.last_name AS name", "r.rental_date, f.title"}
	q := Xql.SELECT(cols...).FROM("rental r")
	q = q.JOIN("inventory i").ON("i.inventory_id = r.inventory_id")
	q = q.JOIN("film f").ON("f.film_id = i.film_id")
	q = q.JOIN("customer c").ON("c.customer_id = r.customer_id")
	return q
}
