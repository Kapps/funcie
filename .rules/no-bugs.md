Make sure our code follows these best practices, UNLESS there's a comment explaining why it's okay to break the rule.

1. Avoid typos.
2. Don't have like, really obvious bugs.
3. Don't store secrets in code.
4. Follow reasonable conventions of the language we're programming in. No need to be too strict.
5. Avoid dangerous stuff, like things that could lead to template injection, SQL injection, broken access control, or really anything that would show up as a CVE somewhere.
6. Avoid situations that could block further requests from coming in. Any request should be processed concurrently and not prevent new ones being processed.
