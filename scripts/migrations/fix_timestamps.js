//used to fix timestamp fields stored in an incorrect format
var result = db.users.updateMany(
    { "created_at.$date": { $exists: true } },
    [{
        $set: {
            created_at: { $toDate: "$created_at.$date" },
            updated_at: { $toDate: "$updated_at.$date" },
            last_login: { 
                $toDate: { 
                    $cond: {
                        if: { $ne: ["$last_login", null] },
                        then: "$last_login.$date",
                        else: new Date(0)
                    }
                }
            },
            last_activity: { 
                $toDate: { 
                    $cond: {
                        if: { $type: "$last_activity.$date" },
                        then: {
                            $cond: {
                                if: { $type: "$last_activity.$date.$numberLong" },
                                then: { $toLong: "$last_activity.$date.$numberLong" },
                                else: "$last_activity.$date"
                            }
                        },
                        else: new Date(0)
                    }
                }
            }
        }
    }]
);

printjson(result);
