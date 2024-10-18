#!/bin/bash
benchstat benchmark/bench_simp_impl.txt benchmark/constant_removals.txt benchmark/bench_better_maps.txt benchmark/bench_maps_impl.txt benchmark/bench_trees_impl.txt > benchmark/compare_benches.txt
