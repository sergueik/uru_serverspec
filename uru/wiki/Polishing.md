## Overview

Uru is a micro-kernel that provides a core set of functionality. In addition to an administrative interface for managing a gaggle-o-rubies, uru is laser focused on
the following tasks:

1. Listing registered rubies available for use.
2. Activating and using any registered ruby.
3. Executing arbitrary `gem` commands against all registered rubies.
4. Executing arbitrary ruby snippets or scripts against all registered rubies.

Effectively, this means that uru's core command set is open to refurbishments and
refinements, but closed to additions. This stance will be relaxed for important additions
to the admin interface. Additional features and capabilities must be delivered as plugins
rather than new core commands.

## Wishlist

1. Hard-core bug hunting safaris.
2. Refactorings focused on making the codebase cleaner, more efficient, more multi-platform
   compatible, or more performant.
3. Tests to allow triple backflips without breaking our necks due to regressions.
4. Creating and implementing a plugin system to enable users to extend the core feature set.
5. Usage documentation

If any of these sound interesting, I'd love your help. If you've got wacky, wild-eyed ideas
for making uru better, I'd also love to hear from you.

Jon