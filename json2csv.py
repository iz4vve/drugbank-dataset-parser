import glob
import json
import os
import sys
import urllib
import pandas as pd

import tqdm

USAGE = """Usage: python json2csv.py <json-glob> <output-directory>"""

def main(args):
    if "-h" in args or len(args) != 2:
        print(USAGE)
        sys.exit(1)

    files = glob.glob(args[0])
    os.makedirs(args[1], exist_ok=True)
    for file_path in tqdm.tqdm(files):
        _, dst = os.path.split(file_path)
        pd.DataFrame(
            [json.loads(i) for i in open(file_path).readlines()]
        ).to_csv(
            os.path.join(args[1], dst.replace(".json", ".csv")),
            index=False
        )
    
    links = pd.read_csv(os.path.join(args[1], "external_links.csv"))[["resource", "url"]]
    links["url"] = links["url"].apply(lambda x: urllib.parse.urlsplit(x).netloc)
    links.drop_duplicates(subset=["resource"]).to_csv(os.path.join(args[1], "external_links_resources.csv"), index=False)

    links = pd.read_csv(os.path.join(args[1], "links.csv"))[["title", "url"]]
    links["url"] = links["url"].apply(lambda x: urllib.parse.urlsplit(x).netloc)
    links.drop_duplicates(subset=["title"]).to_csv(os.path.join(args[1], "links_resources.csv"), index=False)

    identifiers = pd.read_csv(os.path.join(args[1], "external_identifiers.csv"))[["resource"]]
    identifiers.drop_duplicates(subset=["resource"]).to_csv(os.path.join(args[1], "external_identifiers_resource.csv"), index=False)

    organisms = pd.read_csv(os.path.join(args[1], "organisms.csv"))[["organism"]]
    organisms["organism"] = organisms.organism.apply(str.strip)
    organisms.drop_duplicates(subset=["organism"]).to_csv(os.path.join(args[1], "organisms_resources.csv"), index=False)

    packagers = pd.read_csv(os.path.join(args[1], "packagers.csv"))[["name", "url"]]
    packagers["url"] = packagers["url"].astype(str).apply(lambda x: urllib.parse.urlsplit(x).netloc)
    packagers.drop_duplicates(subset=["name"]).to_csv(os.path.join(args[1], "packagers_resources.csv"), index=False)

if __name__ == "__main__":
    try:
        ARGS = sys.argv[1:]
        main(ARGS)
    except Exception as exc:
        print(exc)
        print(USAGE)
