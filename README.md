# Updates
There's a few implementations available on different branches. 
Benchmark folder in root of main branch should have all the raw benchmarking outputs and the comparison generated with benchstat.

## Versioning Overview
This table outlines some of common primitive operations' time complexities.
I opted for this rather than the full functions because many of them require
touching on order of O(n) Folder structs just to update them, meaning they all
have a lower bound of O(n) anyway.

The more important factors are therefore the cost to find a node in the tree and
the cost of changing the ADT structure in a MoveFolder call, hence the table columns.

n = number of folders.

m = number of orgs.

|                                                                                              	| driver build avg TC 	| driver build worst TC 	| lookup Folder avg TC 	| lookup Folder worst TC 	| delete Folder avg TC            | delete Folder worst TC            |
|----------------------------------------------------------------------------------------------	|---------------------	|-----------------------	|----------------------	|------------------------	|-------------------------------	|---------------------------------	|
| [optimal, single threaded](https://github.com/jon-atkinson/sc-takehome-2024-25/tree/main)     | O(n log(n))         	| O(n^2)         	        | O(1)                 	| O(n)                   	| O(1)                     	      | O(n)                          	  |
| [no orgs, all maps](https://github.com/jon-atkinson/sc-takehome-2024-25/tree/org_w_maps_impl) | O(n log(n))         	| O(n^2 log(n))         	| O(1)                 	| O(n)                   	| O(log(n))                     	| O(n^2)                          	|
| [map_trees](https://github.com/jon-atkinson/sc-takehome-2024-25/tree/feat_trees_to_maps)     	| O(n log(n))         	| O(n^2 log(n))         	| O(log(n))            	| O(m n)                 	| O(log(n))                     	| O(m n^2)                        	|
| [btrees](https://github.com/jon-atkinson/sc-takehome-2024-25/tree/trees_impl)                	| O(n log(n))         	| O(n log(n))           	| O(log(m) log(n))     	| O(log(m) log(n))       	| O(log(m) log(n))              	| O(log(m) log(n))                	|
| [simple](https://github.com/jon-atkinson/sc-takehome-2024-25/tree/simplified_implementation) 	| O(1)                	| O(1)                  	| O(n)                 	| O(n)                   	| O(n)                          	| O(n)                            	|

## Benchmarking
This branch contains a benchmark directory that contains benchmarking output
and benchstat comparison for all implementations outlined above.

