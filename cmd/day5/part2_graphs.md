# Example 1
* Invalid order: `75,97,47,61,53`
* Corrected order: `97,75,47,61,53`

```mermaid
graph LR
subgraph "Pages: 75, 97, 47, 61, 53"
75
97
47
61
53
75 --> 97
47 --> 97
47 --> 75
61 --> 97
61 --> 47
61 --> 75
53 --> 47
53 --> 75
53 --> 61
53 --> 97
end
```

# Example 2
* Invalid order: `61,13,29`
* Corrected order: `61,29,13`

```mermaid
graph LR
subgraph "Pages: 61, 13, 29"
61
13
29
13 --> 61
13 --> 29
29 --> 61
end
```

# Example 3
* Invalid order: `97,13,75,29,47`
* Corrected order: `97,75,47,29,13`

```mermaid
graph LR
subgraph "Pages: 97, 13, 75, 29, 47"
97
13
75
29
47
47 --> 97
47 --> 75
13 --> 97
13 --> 29
13 --> 47
13 --> 75
75 --> 97
29 --> 75
29 --> 97
29 --> 47
end
```

