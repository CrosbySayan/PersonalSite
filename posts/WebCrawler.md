# Web Crawler

### Description
This assignment is intended to familiarize you with the [[HTTP]] protocol - the most important application level protocol on the Internet today: the Web runs on HTTP, and increasingly other applications use HTTP as well (including [[Bittorrent]], streaming video, Facebook and Twitter's social APIs, etc). 

Your goal in this assignment is to implement a web crawler that gathers data from a fake social networking website that we have set up for you. Educational goals for the project:
- To expose you to the [[HTTP]] 1.1 protocol, which underlies a large (and growing) number of applications and services today.
- To let you see how web pages are structured using [[HTML]] 
- To give you experience implementing a client for a well-specified network protocol.
- To have you understand how web crawlers work, and how they are used to implement popular web services today.

### Definition: What is a Web Crawler?

A web crawler is a piece of software that automatically gathers and traverses documents on the web. For example, lets say you have a crawler and you tell it to start at https://www.wikipedia.com. The software will first download the Wikipedia homepage, then it will parse the HTML and locate all hyperlinks (i.e., anchor tags) embedded in the page. the crawler then downloads all the HTML pages specified by the URLs on the homepage, and parses them looking for more hyperlinks. This process continues until all of the pages on Wikipedia are downloaded and parsed.

Web crawlers are a fundamental component of today’s web. For example, Googlebot is Google’s web crawler. Googlebot is constantly scouring the web, downloading pages in search of new and updated content. All of this data forms the backbone of Google’s search engine infrastructure.

---
## Project Outline: Fakebook

Set up a fake social network called Fakebook. The website consists of the following pages:
- **Homepage:** The facebook homepage displays some welcome text, as well as links to several random Fakebook users' personal profiles.
- **Personal Profiles:** Each Fakebook user has a profile page that includes their name, some basic demographic information, as well as a link to their list of friends
- **Friends List**: Each Fakebook user is friends with one or more other Fakebook users. This page lists the user's friends and has links to their personal profiles

> [!Tip] WARNING: DO NOT TEST YOUR WEB CRAWLERS ON PUBLIC WEBSITES
> 
> Many web server administrators view crawlers as a nuisance, and they get very mad if they see strange crawlers traversing their sites. **Only test your crawler against Fakebook, do not test it against any other websites.**


#### Username and Password

Your username for Fakebook is your **Northeastern username** not including the @. To get your password visit [this](https://at-tools.khoury.northeastern.edu/course-password/index.php) site and enter your NUID with leading zeroes.

#### High-level Reqs
Collect five *secret flags* that have been hidden somewhere on the Fakebook websites. The flags are unique for each student, and the pages that contain the flags will be different for each student. Since you have no idea what pages the secret flags will appear on, and the Fakebook site is very large (tens of thousands of pages), your only option is to write a web crawler that will traverse Fakebook and locate your flags. 

Your web crawler must execute on the command line using the following syntax: 
``$ ./crawler <-s server> <-p port> <username> <password>`` 

The -s and -p arguments are each optional and they represent the server and port your code should crawl, respectively

Defaults: 
Server ``fakebook.khoury.northeastern.edu
Port ``443
``https://<server>:<port>/fakebook/ 
``https://<server>:<port>/accounts/login/?next=/fakebook/

#### Printing
The crawler should print **exactly five lines of output to STDOUT**: the five secret flags discovered during the crawl of Fakebook, each terminate by a \n character. Your web crawler should not print out anything other than those flags. If a crash occurs error messages are allowed.

**Secrete Flags**
They may be hidden on any page on Fakebook, their exact location on each page may be different, and pages may contain multiple flags. Each secret flag is a 64 character long sequences of random alphanumerics. All secret flags will appear in the following format

``<h3 class='secret_flag style="color:red">FLAG: 64-characters-of-random-alphanumerics</h3>`` 
  
#### Sockets, Ports, and [[TLS]]
Fakebook uses [[HTTPS]], which means that the full protocol stack is HTTP 1.1 over TLS over [[TCP]]. Thus, in this project, your web crawler will need to connect to Fakebook using a TCP socket wrapped in TLS.  Similar to the TLS version of the simple client made in the [[Wordle Socket]], where we made non-secure and secure versions of the client.

HTTP uses TCP port 80 and HTTPS uses port 443
Connect to Fakebook on port 443 unless -p is specified

## Implementation Details and Hints

The crawler must implement HTTP/1.1. This means that there are certain HTTP headers like Host that you must include in your requests (i.e., they are required for all HTTP/1.1 requests). We encourage you to implement **Connection: Keep-Alive** (i.e., pipelining) to improve your crawlers performance (and lighten the load on our server), but this is not required, and it is tricky to get correct.

Also encourage students to implement Accept-Encoding: gzip (i.e., compressed HTTP responses), since this will also improve performance for everyone, but this also not required.

Massive speed boosts can be achieved through parallelism. **Limit to 5 requests though**

One of the key difference btween HTTP/1.0 and HTTP/1.1 is that the latter suppors chunked encoding. HTTP/1.1 servers may break up large response into chunks, and it is the client's responsibility to reconstruct the data by combining the chunks. These servers may return chunked response.

### Key components
- HTTP GET - These requests are necessary for downloading HTML pages
- HTTP POST - You will need to implement HTTP POST so that your code can login to Fakebook. As shown above, you will pass a username and password to your crawler on the command line. The crawler will then use these values as parameters in an HTTP POST in order to log-in to Fakebook.
- Cookie Management - Fakebook uses cookies to track wether clients are logged in to the site. If your crawler successfully logs in to Fakebook using an HTTP POST, Fakebook will return a session cookie to your crawler. Your crawler should store this cookie, and submit it along with each HTTP GET request as it crawls Fakebook. If your crawler fails to handle cookies properly, then your software won't crawl.

**[HTTP Status codes](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes)** 
200 - OK: everything okay
302 - Found: This is as an HTTP redirect. Your crawler should try the request again using the new URL given by the server in the _Location_ header.
403 - Forbidden and 404 - Not Found: Our web server may return these codes in order to trip up your crawler. In this case, your crawler should abandon the URL that generated the error code.

503 - Service unavailable: Our web server may **randomly** return this error code to your crawler. In this case, your crawler should re-try the request for the URL until the request is successful.

If you encounter HTTP `500` errors, please contact the course staff. These may be legitimate bugs in our server software and are not intentional.

---

**Track the Frontier:** As your crawler traverses Fakebook it will observe many URLs. Typically, these uncrawled URLs are stored in a queue, stack, or list until the crawler is ready to visit them. These uncrawled URLs are known as the frontier. Refer to [[Depth-First Search]] or [[Breadth-First Search]] for ideads

**Watch Out for Loops**: Your crawler needs to keep track of where it has been, i.e., the URLs that it has already crawled. Obviously, it isn’t efficient to revisit the same pages over and over again. If your crawler does not keep track of where it has been, it will almost certainly enter an infinite loop. For example, if users A and B are friends on Fakebook, then that means A’s page links to B, and B’s page links to A. Unless the crawler is smart, it will ping-pong back and forth going A->B, B->A, A->B, B->A, …, etc. If you find that your crawler does not find all of its flags and keeps running, you are likely stuck in such a loop. The way to approach this is through [[Graphs]] searching and marking

**Only Crawl The Target Domain**: Web pages may include links that point to arbitrary domains (e.g., a link on google.com that points to cnn.com). **Your crawler should only traverse URLs that point to pages on the specified server (e.g., fakebook.khoury.northeastern.edu, if no `-s` is provided)**. For example, it would be valid to crawl _https://fakebook.khoury.northeastern.edu/fakebook/018912/_, but it would not be valid to crawl _https://www.facebook.com/018912/_. Your code should check to make sure that each URL has a valid domain before you attempt to visit it.

### Logging in to Fakebook
In order to write code that can successfully log-in to Fakebook, you will need to reverse engineer the HTML form on the log-in page. **Students should carefully inspect the form’s code, since it may not be as simple as it initially appears.** The key acronym you should be on the lookout for is _CSRF_.

## Resources
[HTTP Made Really Easy](https://www.jmarshall.com/easy/http/) 
