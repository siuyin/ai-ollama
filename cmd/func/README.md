## Function calling with Gemma3:4b

Sample prompt with good results with both Gemma3:4b and llama3.2:3b.
Gemma3:1b missed out the initial call to Bankok time (short context window?):

```
You have access to functions. If you decide to invoke any of the function(s),
you MUST put it in the json format of
{"name": function name, "parameters": dictionary of argument name and its value}

You SHOULD NOT include any other text in the response if you call a function
[
  {
    "name": "get_product_name_by_PID",
    "description": "Finds the name of a product by its Product ID",
    "parameters": {
      "type": "object",
      "properties": {
        "PID": {
          "type": "string"
        }
      },
      "required": [
        "PID"
      ]
    }
  }, {"name":"getTime", "description":"get the current time in a given timezone", "parameters": {"type":"object", "properties": {"timezone":{"type":"string}} } }
]

What is the time in Bangkok? And while browsing the product catalog, I came across a product that piqued my
interest. The product ID is 807ZPKBL9V. Can you help me find the name of this
product? And also get me the time in Singapore.
```

Run with ollama:

```
eg.
OLLAMA_HOST=http://somehost.com:11434 ollama run gemma3:4b
```
Then copy and paste in the above prompt.


