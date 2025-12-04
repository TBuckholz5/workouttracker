class Workout {
  final int _exerciseID;
  final List<WorkoutSet> _sets;

  Workout({required int exerciseID, required List<WorkoutSet> sets})
    : _sets = sets,
      _exerciseID = exerciseID;

  List<WorkoutSet> get sets => _sets;
  int get exerciseID => _exerciseID;

  Map<String, dynamic> toJson() {
    Map<String, dynamic> res = {"exerciseID": _exerciseID, "sets": []};
    for (WorkoutSet set in _sets) {
      res["sets"].add(set.toJson());
    }
    return res;
  }
}

class WorkoutSet {
  final int _reps;
  final double _weight;
  final String _type;
  final int _order;

  WorkoutSet({
    required int reps,
    required double weight,
    required String type,
    required int order,
  }) : _order = order,
       _type = type,
       _weight = weight,
       _reps = reps;

  int get reps => _reps;
  double get weight => _weight;
  String get type => _type;
  int get order => _order;

  Map<String, dynamic> toJson() {
    return {
      "reps": _reps,
      "weight": _weight,
      "set_type": _type,
      "set_order": _order,
    };
  }
}

class Exercise {
  final int? _id;
  final String? _name;
  final String? _description;
  final String? _targetMuscle;

  int? get id => _id;
  String? get name => _name;
  String? get description => _description;
  String? get targetMuscle => _targetMuscle;

  static Exercise fromJson(Map<String, dynamic> json) {
    return Exercise(
      id: json['id'],
      name: json['name'],
      description: json['description'],
      muscleGroup: json['muscleGroup'],
    );
  }

  Map<String, dynamic> toJson() {
    final map = <String, dynamic>{};

    if (_id != null) map["id"] = _id;
    if (_name != null) map["name"] = _name;
    if (_description != null) map["description"] = _description;
    if (_targetMuscle != null) map["targetMuscle"] = _targetMuscle;

    return map;
  }

  Exercise({int? id, String? name, String? description, String? muscleGroup})
    : _id = id,
      _name = name,
      _description = description,
      _targetMuscle = muscleGroup;
}
