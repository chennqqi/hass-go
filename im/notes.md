# IM becomes Shout!

Calendar currently is focused on sensors, how do we mark things we want to push
on Slack or other messaging platforms ?

Do we move state.Instance into sensors.State ?
And we make a shout.State where we can push messages into like
shout.State.PushMessage(UUID, title, description, body, severity.Critical/Error/Warning/Info)

Calendar should have events like Shout.Message{.Critical} where Description contains 
the description and notes contains the body?

Default severity is Warning

