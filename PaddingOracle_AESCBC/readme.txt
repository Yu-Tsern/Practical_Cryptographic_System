Assume that the attacker is sending packets to the server for decryption, and that he receives some notification from the server when the routine outputs an error (i.e., returns 0). How might he execute a padding oracle attack? (Hint: this is trickier than it looks!)

The first thing came to my mind is to exploit the difference of returning time and power consumption. Since there is no way to check MAC without checking paddings, there must be a small difference in computing time or power consumption between invalid mac and invalid padding.
