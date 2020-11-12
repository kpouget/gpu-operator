#! /bin/bash
exec operator-sdk run local --operator-flags '--zap-encoder=console' |& cut -b35-
