def argparser(parser):
    """Default argument parser for regressions."""
    parser.add_argument(
        "--native",
        action="store_true",
        help="run tests without docker-compose, require only working kubectl + python",
        default=False,
    )
    parser.add_argument(
        "--keeper-type",
        type=str,
        help="type of keeper to use for tests",
        choices=["zookeeper", "clickhouse-keeper"],
        default="zookeeper",
    )
    parser.add_argument(
        "--test-part",
        type=str,
        help="test part to run (part1, part2, part3, part4, no_parallel)",
        choices=["part1", "part2", "part3", "part4", "no_parallel"],
        default="no_parallel",
    )
