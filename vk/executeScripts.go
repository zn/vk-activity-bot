package vk

const GetSubscribersScript string = `
		var offset = parseInt(Args.offset);
		if(offset == null) {
			offset = 0;
		}
		var group_id = Args.group_id;
		var calls = 0;
		var users = [];
		
		while(calls < 25) {
			var received = [API.groups.getMembers({"group_id":group_id,"offset":offset})]@.items[0];
			if(received.length == 0){
				return {"done":true, "offset":offset, "items":users};
			}
			offset = offset + received.length;
			users = users + received;
		
			calls = calls + 1;
		}
		
		return {"done":false,"offset":offset, "items":users};`
