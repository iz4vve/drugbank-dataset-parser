import json
import os
import sys
import pandas as pd

import glob
import tqdm

USAGE = """Usage: python json2csv.py <json-glob> <output-directory>"""

def main(args):
    if "-h" in args or len(args) != 2:
        print(USAGE)
        sys.exit(1)
    
    files = glob.glob(args[0])
    os.makedirs(args[1], exist_ok=True)
    for f in tqdm.tqdm(files):
        directory, dst = os.path.split(f)
        pd.DataFrame(
            [json.loads(i) for i in open(f).readlines()]
        ).to_csv(
            os.path.join(args[1], dst.replace(".json", ".csv")),
            index=False
        )

if __name__ == "__main__":
    try:
        args = sys.argv[1:]
        main(args)
    except Exception as exc:
        print(exc)
        print(USAGE)
