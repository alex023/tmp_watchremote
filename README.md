# tmp_watchremote
test for protoactor-go remote.Watch

# run step
1. start server,client1,client2
   watching:
   - client1 will stop agent a
   - client2 will stop agent c
   - server can not catch terminated event
2. restart server,client1,client2
   - stop client1 or client2
   - server can not catch terminated event
