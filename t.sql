--  some comment wHere, another comment

	select tn.X, sum(tn.Y)  as y ,min(tn.Z) as z from tableName tn 
	left join anotherT t ON tn.X = t.Y Where tn.A = 2
	group                    by tn.X
	order






						by x
	option(recompile, force order)



		;with x as (select * from tableName where x = 10 and y = 10) select * from x

