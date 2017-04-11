# Change Log

## [v1.4.0](https://github.com/nubisproject/nubis-bastionsshkey/tree/v1.4.0) (2017-03-22)
[Full Changelog](https://github.com/nubisproject/nubis-bastionsshkey/compare/v1.3.0...v1.4.0)

**Closed issues:**

- Tag v1.4.0 release [\#59](https://github.com/nubisproject/nubis-bastionsshkey/issues/59)

**Merged pull requests:**

- Merge v1.4.0 release into develop. \[skip ci\] [\#61](https://github.com/nubisproject/nubis-bastionsshkey/pull/61) ([tinnightcap](https://github.com/tinnightcap))
- Update CHANGELOG for v1.4.0 release \[skip ci\] [\#60](https://github.com/nubisproject/nubis-bastionsshkey/pull/60) ([tinnightcap](https://github.com/tinnightcap))

## [v1.3.0](https://github.com/nubisproject/nubis-bastionsshkey/tree/v1.3.0) (2016-12-21)
**Implemented enhancements:**

- Add --version option [\#37](https://github.com/nubisproject/nubis-bastionsshkey/issues/37)
- Update documentation [\#33](https://github.com/nubisproject/nubis-bastionsshkey/issues/33)

**Closed issues:**

- Print out stderr when unicreds command fails [\#55](https://github.com/nubisproject/nubis-bastionsshkey/issues/55)
- Key value store should just use ldap group names instead [\#48](https://github.com/nubisproject/nubis-bastionsshkey/issues/48)
- Ability to create roles [\#27](https://github.com/nubisproject/nubis-bastionsshkey/issues/27)
- If you are using lambda you don't need to pass AWS IAM credentials [\#24](https://github.com/nubisproject/nubis-bastionsshkey/issues/24)
- userCreationPath still hardcoded [\#22](https://github.com/nubisproject/nubis-bastionsshkey/issues/22)
- Weird issues with execType=consul [\#9](https://github.com/nubisproject/nubis-bastionsshkey/issues/9)
- remove hardcoded region [\#8](https://github.com/nubisproject/nubis-bastionsshkey/issues/8)
- Makefile to cross compile package [\#6](https://github.com/nubisproject/nubis-bastionsshkey/issues/6)
- Add support for config file flag [\#4](https://github.com/nubisproject/nubis-bastionsshkey/issues/4)
- Tag v1.3.0 release [\#57](https://github.com/nubisproject/nubis-bastionsshkey/issues/57)

**Merged pull requests:**

- Update CHANGELOG for v1.3.0 release [\#58](https://github.com/nubisproject/nubis-bastionsshkey/pull/58) ([tinnightcap](https://github.com/tinnightcap))
- Printing stderr [\#56](https://github.com/nubisproject/nubis-bastionsshkey/pull/56) ([limed](https://github.com/limed))
- Now properly accepting | delimited list of groups in IAMGroupMapping:… [\#54](https://github.com/nubisproject/nubis-bastionsshkey/pull/54) ([rtucker-mozilla](https://github.com/rtucker-mozilla))
- Added checks and ending execution if usersSet is empty [\#53](https://github.com/nubisproject/nubis-bastionsshkey/pull/53) ([rtucker-mozilla](https://github.com/rtucker-mozilla))
- Use ldap group name for one of the kv store name [\#49](https://github.com/nubisproject/nubis-bastionsshkey/pull/49) ([limed](https://github.com/limed))
- Try and use goveralls.io to report on code coverage [\#46](https://github.com/nubisproject/nubis-bastionsshkey/pull/46) ([gozer](https://github.com/gozer))
- Improve depth of Travis testing [\#44](https://github.com/nubisproject/nubis-bastionsshkey/pull/44) ([gozer](https://github.com/gozer))
- Adding some print lines [\#43](https://github.com/nubisproject/nubis-bastionsshkey/pull/43) ([limed](https://github.com/limed))
- Enable simple Travis-CI integration [\#42](https://github.com/nubisproject/nubis-bastionsshkey/pull/42) ([gozer](https://github.com/gozer))
- Have a way to show what version we are running [\#39](https://github.com/nubisproject/nubis-bastionsshkey/pull/39) ([limed](https://github.com/limed))
- Doc update [\#35](https://github.com/nubisproject/nubis-bastionsshkey/pull/35) ([limed](https://github.com/limed))
- enable email sending [\#31](https://github.com/nubisproject/nubis-bastionsshkey/pull/31) ([limed](https://github.com/limed))
- Some formatting changes [\#30](https://github.com/nubisproject/nubis-bastionsshkey/pull/30) ([limed](https://github.com/limed))
- Refactored adding iam users and addition of tests [\#29](https://github.com/nubisproject/nubis-bastionsshkey/pull/29) ([rtucker-mozilla](https://github.com/rtucker-mozilla))
- Added ability to create and delete roles [\#28](https://github.com/nubisproject/nubis-bastionsshkey/pull/28) ([limed](https://github.com/limed))
- Wrote function to return session.Session object and share it across IAM [\#26](https://github.com/nubisproject/nubis-bastionsshkey/pull/26) ([rtucker-mozilla](https://github.com/rtucker-mozilla))
- Added functionality to remove users from IAM [\#25](https://github.com/nubisproject/nubis-bastionsshkey/pull/25) ([rtucker-mozilla](https://github.com/rtucker-mozilla))
- Added tracking of userPaths and users [\#23](https://github.com/nubisproject/nubis-bastionsshkey/pull/23) ([rtucker-mozilla](https://github.com/rtucker-mozilla))
- Proper code in place to create users in IAM from LDAP [\#20](https://github.com/nubisproject/nubis-bastionsshkey/pull/20) ([rtucker-mozilla](https://github.com/rtucker-mozilla))
- Set port when using consulPort [\#19](https://github.com/nubisproject/nubis-bastionsshkey/pull/19) ([limed](https://github.com/limed))
- Some make file updates / enhancements [\#18](https://github.com/nubisproject/nubis-bastionsshkey/pull/18) ([limed](https://github.com/limed))
- Added ability back to provide ConsulPort with a default [\#17](https://github.com/nubisproject/nubis-bastionsshkey/pull/17) ([rtucker-mozilla](https://github.com/rtucker-mozilla))
- Removed dupe functions [\#16](https://github.com/nubisproject/nubis-bastionsshkey/pull/16) ([rtucker-mozilla](https://github.com/rtucker-mozilla))
- Moved remaining ancillary functions from main.go [\#15](https://github.com/nubisproject/nubis-bastionsshkey/pull/15) ([rtucker-mozilla](https://github.com/rtucker-mozilla))
- Lambda specific add [\#14](https://github.com/nubisproject/nubis-bastionsshkey/pull/14) ([limed](https://github.com/limed))
- Tiny nitpicks [\#13](https://github.com/nubisproject/nubis-bastionsshkey/pull/13) ([limed](https://github.com/limed))
- Missing unicreds param [\#12](https://github.com/nubisproject/nubis-bastionsshkey/pull/12) ([limed](https://github.com/limed))
- Cleanup makefile [\#11](https://github.com/nubisproject/nubis-bastionsshkey/pull/11) ([limed](https://github.com/limed))
- Add missing unicred argument [\#10](https://github.com/nubisproject/nubis-bastionsshkey/pull/10) ([limed](https://github.com/limed))
- Added a Makefile to build binary [\#7](https://github.com/nubisproject/nubis-bastionsshkey/pull/7) ([limed](https://github.com/limed))
- Flattened consul KV to remove serial and doing a compare of sshPublic… [\#5](https://github.com/nubisproject/nubis-bastionsshkey/pull/5) ([rtucker-mozilla](https://github.com/rtucker-mozilla))
- Instead of getting users in a group, just get the enabled users [\#3](https://github.com/nubisproject/nubis-bastionsshkey/pull/3) ([limed](https://github.com/limed))
- Add README.md [\#2](https://github.com/nubisproject/nubis-bastionsshkey/pull/2) ([limed](https://github.com/limed))
- Config.yml-dist verbiage update [\#1](https://github.com/nubisproject/nubis-bastionsshkey/pull/1) ([limed](https://github.com/limed))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*