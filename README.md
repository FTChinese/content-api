## Byline

When generating data for byline, editors should follow the convention to use English comma `,` to separate groups, and use semicolon `;` to separate sub-groups

The original version:

```
James Politi in Washington and Alistair Gray, Richard Henderson and Fan Fei in New York
```

When translated into Chinese, you should write it this way:

```
詹姆斯•波利提,阿利斯泰尔•格雷;理查德•亨德森

华盛顿,纽约报道
```

It will produce the correct Chinese text:
```
詹姆斯•波利提 华盛顿， 阿利斯泰尔•格雷，理查德•亨德森 纽约报道
```

If you separate everything with comma, it will produce
```
詹姆斯•波利提 华盛顿， 阿利斯泰尔•格雷 纽约报道， 理查德•亨德森
```

which is apparently not correct.
