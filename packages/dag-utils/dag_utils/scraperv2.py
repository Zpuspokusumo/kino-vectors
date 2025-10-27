import requests
from bs4 import BeautifulSoup

class MovieInfo:
    def __init__(self, title="", director="", year="", genre="", actors=None, summary=""):
        self.title = title
        self.director = director
        self.year = year
        self.genre = genre
        self.actors = actors if actors else []
        self.summary = summary

    def __repr__(self):
        return f"MovieInfo(title={self.title}, director={self.director}, year={self.year}, genre={self.genre}, actors={self.actors}, summary={self.summary})"

def logtofile(title: str, msg: str):
    with open("sqr2.log", "a", encoding="utf-8") as f:
        f.write(f"{title}: {msg}\n")

from bs4 import BeautifulSoup

def extract_plot(soup):
    # find the heading for Plot
    plot_header = soup.find("h2", id="Plot")
    if not plot_header:
        return None  # or empty string if you prefer

    paragraphs = []
    # iterate over the siblings until we hit another h2
    for sibling in plot_header.find_all_next():
        # stop when we reach the next main section
        if sibling.name == "h2":
            break
        if sibling.name == "p":
            paragraphs.append(sibling.get_text(" ", strip=True))

    return "\n\n".join(paragraphs)


def fetch_movie_info(title: str) -> MovieInfo:
    headers = {
        "User-Agent": "kino-vector/0.1 (contact: your_email@example.com)"
    }

    # Fetch whole page HTML so we can parse the infobox
    url = "https://en.wikipedia.org/w/api.php"
    params = {
        "action": "parse",
        "page": title,
        "format": "json",
        "prop": "text",
    }
    r = requests.get(url, params=params, headers=headers)
    r.raise_for_status()
    html = r.json()["parse"]["text"]["*"]

    soup = BeautifulSoup(html, "html.parser")

    # Try to find infobox
    infobox = soup.find("table", {"class": "infobox vevent"})
    info = MovieInfo()

    if infobox:
        for row in infobox.find_all("tr"):
            header = row.find("th")
            value = row.find("td")
            if not header or not value:
                continue

            key = header.get_text(strip=True).lower()
            val = value.get_text(" ", strip=True)

            if "directed by" in key:
                info.director = val
            elif "release date" in key or "released" in key:
                info.year = val.split()[-1]  # naive: grab last token as year
            elif "starring" in key:
                info.actors = [a.strip() for a in val.split(",")]
            elif "genre" in key:  # not always present in infobox
                info.genre = val
            elif "title" in key:
                info.title = val

    # Fallback title if not found in infobox
    if not info.title:
        info.title = title

    # Fetch summary extract
    # params_summary = {
    #     "action": "query",
    #     "prop": "extracts",
    #     "exintro": True,
    #     "explaintext": True,
    #     "titles": title,
    #     "format": "json",
    # }
    # r2 = requests.get(url, params=params_summary, headers=headers)
    # r2.raise_for_status()
    # pages = r2.json()["query"]["pages"]
    # page = next(iter(pages.values()))
    
    # paragraphs = [p.get_text(" ", strip=True) for p in soup.find_all("p")]
    # clean_text = "\n\n".join(paragraphs)
    # info.summary = clean_text
    info.summary = extract_plot(soup)


    return info


if __name__ == "__main__":
    movie = fetch_movie_info("Inception")
    print(movie)
